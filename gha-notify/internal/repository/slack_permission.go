package repository

import "context"

// SlackPermissionGetter is an interface for getting slack permission.
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

type SlackPermissionAllower interface {
	AllowSlackPermission(ctx context.Context, input *AllowSlackPermissionInput) (*AllowSlackPermissionOutput, error)
}

type AllowSlackPermissionInput struct {
	TeamID    string
	ChannelID string
	Repos     []string
}

type AllowSlackPermissionOutput struct {
}
