package handler

import (
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

	signingSecret, err := h.cfg.GetSlackSigningSecret(ctx, &repository.GetSlackSigningSecretInput{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret.SlackSigningSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r = shallowCopy(r)
	r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := verifier.Ensure(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.cfg.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(h.cfg.SlashFunctionName),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payload,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
