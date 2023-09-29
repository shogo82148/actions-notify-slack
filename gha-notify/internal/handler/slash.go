package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
	"github.com/slack-go/slack"
)

type SlashHandler struct {
	cfg *SlashHandlerConfig
}

type SlashHandlerConfig struct {
	repository.SlackSigningSecretGetter
	service.LambdaInvoker
	SlashFunctionName string
}

func NewSlashHandler(cfg *SlashHandlerConfig) (*SlashHandler, error) {
	return &SlashHandler{
		cfg: cfg,
	}, nil
}

func shallowCopy[T any](v *T) *T {
	w := *v
	return &w
}

func (h *SlashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	s, err := h.parseRequest(ctx, r)
	if err != nil {
		handleError(ctx, w, err)
		return
	}

	if err := h.handleCommand(ctx, s); err != nil {
		handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SlashHandler) parseRequest(ctx context.Context, r *http.Request) (*slack.SlashCommand, error) {
	signingSecret, err := h.cfg.GetSlackSigningSecret(ctx, &repository.GetSlackSigningSecretInput{})
	if err != nil {
		return nil, err
	}

	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret.SlackSigningSecret)
	if err != nil {
		return nil, newValidationError(err)
	}

	r = shallowCopy(r)
	r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		return nil, err
	}
	if err := verifier.Ensure(); err != nil {
		return nil, newValidationError(err)
	}
	return &s, nil
}

func (h *SlashHandler) handleCommand(ctx context.Context, s *slack.SlashCommand) error {
	payload, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = h.cfg.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(h.cfg.SlashFunctionName),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payload,
	})
	if err != nil {
		return err
	}
	return nil
}
