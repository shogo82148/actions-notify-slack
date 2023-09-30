package repository

import "context"

type SlackPermissionGetter interface {
	GetSlackPermission(ctx context.Context, input *GetSlackPermissionInput) (*GetSlackPermissionOutput, error)
}

type GetSlackPermissionInput struct {
	TeamID    string
	ChannelID string
}

type GetSlackPermissionOutput struct {
	TeamID    string
	ChannelID string
	Repos     []string
}
