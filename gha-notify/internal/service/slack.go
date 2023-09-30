package service

import "context"

// SlackOAuthV2ResponseGetter is an interface for slack.GetOAuthV2Response.
type OAuthV2ResponseGetter interface {
	GetOAuthV2Response(ctx context.Context, input *GetOAuthV2ResponseInput) (*GetOAuthV2ResponseOutput, error)
}

type GetOAuthV2ResponseInput struct {
	ClientID     string
	ClientSecret string
	Code         string
	RedirectURI  string
}

type GetOAuthV2ResponseOutput struct {
	AccessToken  string
	TokenType    string
	Scope        string
	BotUserID    string
	TeamID       string
	TeamName     string
	RefreshToken string
	ExpiresIn    int
}

// OAuthV2ResponseRefresher is an interface for slack.RefreshOAuthV2Token.
type OAuthV2ResponseRefresher interface {
	RefreshOAuthV2Response(ctx context.Context, input *RefreshOAuthV2ResponseInput) (*RefreshOAuthV2ResponseOutput, error)
}

type RefreshOAuthV2ResponseInput struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
}

type RefreshOAuthV2ResponseOutput struct {
	AccessToken  string
	TokenType    string
	Scope        string
	BotUserID    string
	TeamID       string
	RefreshToken string
	ExpiresIn    int
}

// SlackWebhookPoster is an interface for slack.Webhook.Post.
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

// SlackMessagePoster is an interface for slack.PostMessage.
type SlackMessagePoster interface {
	PostSlackMessage(ctx context.Context, input *PostSlackMessageInput) (*PostSlackMessageOutput, error)
}

type PostSlackMessageInput struct {
	Token   string
	Message map[string]any
}

type PostSlackMessageOutput struct {
}
