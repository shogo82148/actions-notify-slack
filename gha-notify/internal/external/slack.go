package external

import (
	"context"
	"net/http"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
	"github.com/slack-go/slack"
)

type Slack struct {
	cfg *SlackConfig
}

type SlackConfig struct {
	HTTPClient *http.Client
}

func NewSlack(cfg *SlackConfig) (*Slack, error) {
	return &Slack{cfg: cfg}, nil
}

func (s *Slack) PostSlackWebhook(ctx context.Context, input *service.PostSlackWebhookInput) (*service.PostSlackWebhookOutput, error) {
	err := slack.PostWebhookCustomHTTPContext(ctx, input.WebhookURL, s.cfg.HTTPClient, &slack.WebhookMessage{
		Text:         input.Text,
		ResponseType: input.ResponseType,
	})
	return nil, err
}
