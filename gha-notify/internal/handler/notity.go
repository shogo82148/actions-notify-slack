package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
)

type NotifyHandler struct {
	cfg *NotifyHandlerConfig
}

type NotifyHandlerConfig struct {
	service.OAuthV2ResponseRefresher
	service.SlackMessagePoster
	repository.SlackClientIDGetter
	repository.SlackClientSecretGetter
	repository.SlackAccessTokenGetter
	repository.SlackAccessTokenPutter
}

func NewNotifyHandler(cfg *NotifyHandlerConfig) (*NotifyHandler, error) {
	return &NotifyHandler{
		cfg: cfg,
	}, nil
}

func (h *NotifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := h.handle(ctx, r); err != nil {
		handleError(ctx, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *NotifyHandler) handle(ctx context.Context, r *http.Request) error {
	// TODO: authorize the request

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		return newValidationError(err)
	}

	teamID, ok := v["team"].(string)
	if !ok {
		return newValidationError(errors.New("handler: required key team is not found"))
	}
	delete(v, "team")

	token, err := h.getAccessToke(ctx, time.Now(), teamID)
	if err != nil {
		return err
	}

	// send a message
	h.cfg.PostSlackMessage(ctx, &service.PostSlackMessageInput{
		Token:   token,
		Message: v,
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *NotifyHandler) getAccessToke(ctx context.Context, now time.Time, teamID string) (string, error) {
	out, err := h.cfg.GetSlackAccessToken(ctx, &repository.GetSlackAccessTokenInput{
		TeamID: teamID,
	})
	if err != nil {
		return "", err
	}
	if out.ExpiresAt.Compare(now) > 0 {
		// the access token is still valid
		return out.AccessToken, nil
	}

	// need to refresh the access token
	clientID, err := h.cfg.GetSlackClientID(ctx, &repository.GetSlackClientIDInput{})
	if err != nil {
		return "", err
	}
	clientSecret, err := h.cfg.GetSlackClientSecret(ctx, &repository.GetSlackClientSecretInput{})
	if err != nil {
		return "", err
	}
	refreshed, err := h.cfg.RefreshOAuthV2Response(ctx, &service.RefreshOAuthV2ResponseInput{
		ClientID:     clientID.SlackClientID,
		ClientSecret: clientSecret.SlackClientSecret,
		RefreshToken: out.RefreshToken,
	})
	if err != nil {
		return "", err
	}

	// save the refreshed access token
	expiresAt := now.Add(time.Duration(refreshed.ExpiresIn) * time.Second)
	_, err = h.cfg.PutSlackAccessToken(ctx, &repository.PutSlackAccessTokenInput{
		TeamID:       refreshed.TeamID,
		BotUserID:    refreshed.BotUserID,
		AccessToken:  refreshed.AccessToken,
		Scope:        refreshed.Scope,
		RefreshToken: refreshed.RefreshToken,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		return "", err
	}
	return refreshed.AccessToken, nil
}
