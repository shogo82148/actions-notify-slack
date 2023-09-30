package repository

import "context"

// SessionPutter is an interface for putting session.
type SessionPutter interface {
	PutSession(ctx context.Context, input *PutSessionInput) (*PutSessionOutput, error)
}

type PutSessionInput struct {
	SessionID string
	State     string // state for OAuth2
	TeamID    string // team ID for Slack
	TeamName  string // team name for Slack
}

type PutSessionOutput struct {
}

// SessionGetter is an interface for getting session.
type SessionGetter interface {
	GetSession(ctx context.Context, input *GetSessionInput) (*GetSessionOutput, error)
}

type GetSessionInput struct {
	SessionID string
}

type GetSessionOutput struct {
	SessionID string
	State     string // state for OAuth2
	TeamID    string // team ID for Slack
	TeamName  string // team name for Slack
}
