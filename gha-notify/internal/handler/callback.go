package handler

import (
	"context"
	"errors"
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
	header, err := h.handle(ctx, r)
	if err != nil {
		handleError(ctx, w, err)
		return
	}

	w.Header().Add("Set-Cookie", header)
	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *CallbackHandler) handle(ctx context.Context, r *http.Request) (string, error) {
	now := time.Now()
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// get the session
	s, err := getSession(r, h.cfg.SessionGetter)
	if err != nil {
		return "", err
	}
	if s.SessionID == "" || s.State == "" || s.State != state {
		return "", newValidationError(errors.New("handler: invalid session"))
	}

	clientID, err := h.cfg.GetSlackClientID(ctx, &repository.GetSlackClientIDInput{})
	if err != nil {
		return "", err
	}
	clientSecret, err := h.cfg.GetSlackClientSecret(ctx, &repository.GetSlackClientSecretInput{})
	if err != nil {
		return "", err
	}

	resp, err := h.cfg.GetOAuthV2Response(ctx, &service.GetOAuthV2ResponseInput{
		ClientID:     clientID.SlackClientID,
		ClientSecret: clientSecret.SlackClientSecret,
		Code:         code,
	})
	if err != nil {
		return "", err
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
		return "", err
	}

	// save the session
	s = &session{
		SessionID: newSessionID(),
		State:     newState(),
		TeamID:    resp.TeamID,
		TeamName:  resp.TeamName,
	}
	header, err := putSession(r, h.cfg.SessionPutter, s)
	if err != nil {
		return "", err
	}
	return header, nil
}
