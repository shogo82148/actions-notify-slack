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

type SlackSigningSecretGetter interface {
	GetSlackSigningSecret(ctx context.Context, input *GetSlackSigningSecretInput) (*GetSlackSigningSecretOutput, error)
}

type GetSlackSigningSecretInput struct{}

type GetSlackSigningSecretOutput struct {
	SlackSigningSecret string
}
