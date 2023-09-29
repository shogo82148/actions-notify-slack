package handler

import (
	"net/http"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	"github.com/slack-go/slack"
)

type NotifyHandler struct {
	cfg *NotifyHandlerConfig
}

type NotifyHandlerConfig struct {
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

	// TODO: get the parameters from the request
	teamID := "T3G1HAY66"
	channelD := "C3GMGG162"

	// get the access token
	out, err := h.cfg.GetSlackAccessToken(ctx, &repository.GetSlackAccessTokenInput{
		TeamID: teamID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := out.AccessToken

	// send a message
	api := slack.New(token)
	_, _, err = api.PostMessageContext(ctx, channelD, slack.MsgOptionText("Hello, World!", false))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
