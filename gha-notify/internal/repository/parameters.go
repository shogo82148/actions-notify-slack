package repository

import "context"

type SlackClientIDGetter interface {
	GetSlackClientID(ctx context.Context, input *GetSlackClientIDInput) (*GetSlackClientIDOutput, error)
}

type GetSlackClientIDInput struct{}

type GetSlackClientIDOutput struct {
	SlackClientID string
}

type SlackClientSecretGetter interface {
	GetSlackClientSecret(ctx context.Context, input *GetSlackClientSecretInput) (*GetSlackClientSecretOutput, error)
}

type GetSlackClientSecretInput struct{}

type GetSlackClientSecretOutput struct {
	SlackClientSecret string
}
