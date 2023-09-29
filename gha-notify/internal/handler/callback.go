package handler

import (
	"net/http"
	"time"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	"github.com/slack-go/slack"
)

type CallbackHandler struct {
	cfg *CallbackHandlerConfig
}

type CallbackHandlerConfig struct {
	repository.SlackClientIDGetter
	repository.SlackClientSecretGetter
	repository.SlackAccessTokenPutter
}

func NewCallbackHandler(cfg *CallbackHandlerConfig) (*CallbackHandler, error) {
	return &CallbackHandler{cfg: cfg}, nil
}

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	now := time.Now()

	code := r.URL.Query().Get("code")
	clientID, err := h.cfg.GetSlackClientID(ctx, &repository.GetSlackClientIDInput{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	clientSecret, err := h.cfg.GetSlackClientSecret(ctx, &repository.GetSlackClientSecretInput{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := slack.GetOAuthV2ResponseContext(ctx, http.DefaultClient, clientID.SlackClientID, clientSecret.SlackClientSecret, code, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expiresAt := now.Add(time.Duration(resp.ExpiresIn) * time.Second)
	_, err = h.cfg.PutSlackAccessToken(ctx, &repository.PutSlackAccessTokenInput{
		TeamID:       resp.Team.ID,
		BotUserID:    resp.BotUserID,
		AccessToken:  resp.AccessToken,
		Scope:        resp.Scope,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
