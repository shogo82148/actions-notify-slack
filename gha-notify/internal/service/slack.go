package service

import "context"

type SlackWebhookPoster interface {
	PostSlackWebhook(ctx context.Context, input *PostSlackWebhookInput) (*PostSlackWebhookOutput, error)
}

type PostSlackWebhookInput struct {
	WebhookURL   string
	Text         string
	ResponseType string
}

type PostSlackWebhookOutput struct {
}
