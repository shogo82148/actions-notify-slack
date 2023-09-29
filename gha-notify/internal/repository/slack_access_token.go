package repository

import (
	"context"
	"time"
)

type SlackAccessTokenPutter interface {
	PutSlackAccessToken(ctx context.Context, input *PutSlackAccessTokenInput) (*PutSlackAccessTokenOutput, error)
}

type PutSlackAccessTokenInput struct {
	TeamID       string
	BotUserID    string
	AccessToken  string
	Scope        string
	RefreshToken string
	ExpiresAt    time.Time
}

type PutSlackAccessTokenOutput struct {
}

type SlackAccessTokenGetter interface {
	GetSlackAccessToken(ctx context.Context, input *GetSlackAccessTokenInput) (*GetSlackAccessTokenOutput, error)
}

type GetSlackAccessTokenInput struct {
	TeamID string
}

type GetSlackAccessTokenOutput struct {
	TeamID       string
	BotUserID    string
	AccessToken  string
	Scope        string
	RefreshToken string
	ExpiresAt    time.Time
}
