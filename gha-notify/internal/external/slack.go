package external

import (
	"context"
	"net/http"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
	"github.com/slack-go/slack"
)

var _ service.SlackWebhookPoster = (*Slack)(nil)
var _ service.OAuthV2ResponseGetter = (*Slack)(nil)
var _ service.OAuthV2ResponseRefresher = (*Slack)(nil)

type Slack struct {
	cfg *SlackConfig
}

type SlackConfig struct {
	HTTPClient *http.Client
}

func NewSlack(cfg *SlackConfig) (*Slack, error) {
	return &Slack{cfg: cfg}, nil
}

// GetOAuthV2Response gets an OAuthV2Response from Slack.
func (s *Slack) GetOAuthV2Response(ctx context.Context, input *service.GetOAuthV2ResponseInput) (*service.GetOAuthV2ResponseOutput, error) {
	resp, err := slack.GetOAuthV2ResponseContext(ctx, s.cfg.HTTPClient, input.ClientID, input.ClientSecret, input.Code, input.RedirectURI)
	if err != nil {
		return nil, err
	}
	return &service.GetOAuthV2ResponseOutput{
		AccessToken:  resp.AccessToken,
		TokenType:    resp.TokenType,
		Scope:        resp.Scope,
		BotUserID:    resp.BotUserID,
		TeamID:       resp.Team.ID,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}, nil
}

func (s *Slack) RefreshOAuthV2Response(ctx context.Context, input *service.RefreshOAuthV2ResponseInput) (*service.RefreshOAuthV2ResponseOutput, error) {
	resp, err := slack.RefreshOAuthV2TokenContext(ctx, s.cfg.HTTPClient, input.ClientID, input.ClientSecret, input.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &service.RefreshOAuthV2ResponseOutput{
		AccessToken:  resp.AccessToken,
		TokenType:    resp.TokenType,
		Scope:        resp.Scope,
		BotUserID:    resp.BotUserID,
		TeamID:       resp.Team.ID,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}, nil
}

// PostSlackWebhook posts a message to Slack using Incoming Webhooks.
func (s *Slack) PostSlackWebhook(ctx context.Context, input *service.PostSlackWebhookInput) (*service.PostSlackWebhookOutput, error) {
	err := slack.PostWebhookCustomHTTPContext(ctx, input.WebhookURL, s.cfg.HTTPClient, &slack.WebhookMessage{
		Text:         input.Text,
		ResponseType: input.ResponseType,
	})
	return nil, err
}
