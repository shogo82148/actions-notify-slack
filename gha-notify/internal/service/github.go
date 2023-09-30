package service

import (
	"context"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/model"
)

type GitHubIDTokenParser interface {
	ParseGitHubIDToken(ctx context.Context, input *ParseGitHubIDTokenInput) (*ParseGitHubIDTokenOutput, error)
}

type ParseGitHubIDTokenInput struct {
	IDToken string
}

type ParseGitHubIDTokenOutput struct {
	Claims *model.ActionsIDToken
}
