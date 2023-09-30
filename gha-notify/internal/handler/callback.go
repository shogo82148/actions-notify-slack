package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
)

type CallbackHandler struct {
	cfg *CallbackHandlerConfig
}

type CallbackHandlerConfig struct {
	service.OAuthV2ResponseGetter
	repository.SlackClientIDGetter
	repository.SlackClientSecretGetter
	repository.SlackAccessTokenPutter
	repository.SessionGetter
	repository.SessionPutter
}

func NewCallbackHandler(cfg *CallbackHandlerConfig) (*CallbackHandler, error) {
	return &CallbackHandler{cfg: cfg}, nil
}

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := r.URL.Query().Get("code")
	if err := h.handle(ctx, code); err != nil {
		handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CallbackHandler) handle(ctx context.Context, code string) error {
	now := time.Now()

	clientID, err := h.cfg.GetSlackClientID(ctx, &repository.GetSlackClientIDInput{})
	if err != nil {
		return err
	}
	clientSecret, err := h.cfg.GetSlackClientSecret(ctx, &repository.GetSlackClientSecretInput{})
	if err != nil {
		return err
	}

	resp, err := h.cfg.GetOAuthV2Response(ctx, &service.GetOAuthV2ResponseInput{
		ClientID:     clientID.SlackClientID,
		ClientSecret: clientSecret.SlackClientSecret,
		Code:         code,
	})
	if err != nil {
		return err
	}

	expiresAt := now.Add(time.Duration(resp.ExpiresIn) * time.Second)
	_, err = h.cfg.PutSlackAccessToken(ctx, &repository.PutSlackAccessTokenInput{
		TeamID:       resp.TeamID,
		BotUserID:    resp.BotUserID,
		AccessToken:  resp.AccessToken,
		Scope:        resp.Scope,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		return err
	}
	return nil
}
