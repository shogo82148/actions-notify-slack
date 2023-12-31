package external

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/model"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
	_ "github.com/shogo82148/goat/jwa/rs" // for RS256
	"github.com/shogo82148/goat/jws"
	"github.com/shogo82148/goat/jwt"
	"github.com/shogo82148/goat/oidc"
	"github.com/shogo82148/goat/sig"
)

const (
	// The value of User-Agent header
	githubUserAgent = "actions-notify-slack/1.0"

	// issuer of JWT tokens
	oidcIssuer = "https://token.actions.githubusercontent.com"
)

type GitHub struct {
	cfg *GitHubConfig

	// configure of OpenID Connect
	oidcClient *oidc.Client
}

type GitHubConfig struct {
	HTTPClient *http.Client
}

func NewGitHub(cfg *GitHubConfig) (*GitHub, error) {
	oidcClient, err := oidc.NewClient(&oidc.ClientConfig{
		Doer:      cfg.HTTPClient,
		Issuer:    oidcIssuer,
		UserAgent: githubUserAgent,
	})
	if err != nil {
		return nil, err
	}
	return &GitHub{
		cfg:        cfg,
		oidcClient: oidcClient,
	}, nil
}

func (g *GitHub) ParseGitHubIDToken(ctx context.Context, input *service.ParseGitHubIDTokenInput) (*service.ParseGitHubIDTokenOutput, error) {
	// get the JSON Web Key Set
	set, err := g.oidcClient.GetJWKS(ctx)
	if err != nil {
		return nil, err
	}

	// decode the ID token
	token, err := jwt.Parse([]byte(input.IDToken), jwt.FindKeyFunc(func(header *jws.Header) (key sig.SigningKey, err error) {
		jwk, ok := set.Find(header.KeyID())
		if !ok {
			return nil, fmt.Errorf("github: kid %s is not found", header.KeyID())
		}
		if jwk.Algorithm() != "" && header.Algorithm().KeyAlgorithm() != jwk.Algorithm() {
			return nil, fmt.Errorf("github: alg parameter mismatch")
		}
		key = header.Algorithm().New().NewSigningKey(jwk)
		return
	}))
	if err != nil {
		return nil, fmt.Errorf("github: failed to parse id token: %w", err)
	}

	// verify the ID token
	if token.Claims.Issuer != oidcIssuer {
		return nil, fmt.Errorf("github: invalid issuer %q", token.Claims.Issuer)
	}
	if !slices.Contains(token.Claims.Audience, input.Audience) {
		return nil, errors.New("github: invalid audience")
	}

	var claims model.ActionsIDToken
	if err := token.Claims.DecodeCustom(&claims); err != nil {
		return nil, fmt.Errorf("github: failed to decode claims: %w", err)
	}
	claims.Claims = token.Claims
	if !strings.HasPrefix(claims.Subject, fmt.Sprintf("repo:%s:", claims.Repository)) {
		return nil, fmt.Errorf("github: invalid subject %q", claims.Subject)
	}

	return &service.ParseGitHubIDTokenOutput{
		Claims: &claims,
	}, nil
}
