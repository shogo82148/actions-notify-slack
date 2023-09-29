package database

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
	"github.com/shogo82148/memoize"
)

var _ repository.SlackClientIDGetter = (*Parameters)(nil)

type Parameters struct {
	g   memoize.Group[string, string]
	cfg *ParametersConfig
}

type ParametersConfig struct {
	service.SSMParameterGetter
}

func NewParameters(cfg *ParametersConfig) (*Parameters, error) {
	return &Parameters{cfg: cfg}, nil
}

func (p *Parameters) getParameter(ctx context.Context, name string) (v string, expiresAt time.Time, err error) {
	param, err := p.cfg.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt = time.Now().Add(10 * time.Minute)
	expiresAt = expiresAt.Round(0) // drop monotonic clock
	return aws.ToString(param.Parameter.Value), expiresAt, nil
}

// GetSlackClientID returns the client ID of the Slack App.
func (p *Parameters) GetSlackClientID(ctx context.Context, _ *repository.GetSlackClientIDInput) (*repository.GetSlackClientIDOutput, error) {
	v, _, err := p.g.Do(ctx, "/slack/client_id", p.getParameter)
	if err != nil {
		return nil, err
	}
	return &repository.GetSlackClientIDOutput{
		SlackClientID: v,
	}, nil
}

// GetSlackClientSecret returns the client secret of the Slack App.
func (p *Parameters) GetSlackClientSecret(ctx context.Context, _ *repository.GetSlackClientSecretInput) (*repository.GetSlackClientSecretOutput, error) {
	v, _, err := p.g.Do(ctx, "/slack/client_secret", p.getParameter)
	if err != nil {
		return nil, err
	}
	return &repository.GetSlackClientSecretOutput{
		SlackClientSecret: v,
	}, nil
}
