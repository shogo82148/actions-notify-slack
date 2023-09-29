package mock

import (
	"context"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
)

var _ repository.SlackSigningSecretGetter = SlackSigningSecretGetterFunc(nil)

type SlackSigningSecretGetterFunc func(ctx context.Context, input *repository.GetSlackSigningSecretInput) (*repository.GetSlackSigningSecretOutput, error)

func (f SlackSigningSecretGetterFunc) GetSlackSigningSecret(ctx context.Context, input *repository.GetSlackSigningSecretInput) (*repository.GetSlackSigningSecretOutput, error) {
	return f(ctx, input)
}
