package handler

import (
	"context"
	"strings"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
	"github.com/slack-go/slack"
)

type SlashCommandHandler struct {
	cfg *SlashCommandHandlerConfig
}

type SlashCommandHandlerConfig struct {
	service.SlackWebhookPoster
}

func NewSlashCommandHandler(cfg *SlashCommandHandlerConfig) (*SlashCommandHandler, error) {
	return &SlashCommandHandler{
		cfg: cfg,
	}, nil
}

func (h *SlashCommandHandler) Handle(ctx context.Context, slash *slack.SlashCommand) (string, error) {
	text := strings.TrimSpace(slash.Text)
	if text == "" || text == "help" {
		return h.handleHelp(ctx, slash)
	}
	return "", nil
}

func (h *SlashCommandHandler) handleHelp(ctx context.Context, slash *slack.SlashCommand) (string, error) {
	h.cfg.PostSlackWebhook(ctx, &service.PostSlackWebhookInput{
		WebhookURL: slash.ResponseURL,
		Text:       "TODO: help message",
	})
	return "", nil
}
