package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
	"github.com/slack-go/slack"
)

type SlashCommandHandler struct {
	cfg *SlashCommandHandlerConfig
}

type SlashCommandHandlerConfig struct {
	service.SlackWebhookPoster
	repository.SlackPermissionGetter
}

func NewSlashCommandHandler(cfg *SlashCommandHandlerConfig) (*SlashCommandHandler, error) {
	return &SlashCommandHandler{
		cfg: cfg,
	}, nil
}

func (h *SlashCommandHandler) Handle(ctx context.Context, slash *slack.SlashCommand) (string, error) {
	if err := h.handle(ctx, slash); err != nil {
		slog.ErrorContext(ctx, "slash command error", slog.String("error", err.Error()))
		h.cfg.PostSlackWebhook(ctx, &service.PostSlackWebhookInput{
			WebhookURL: slash.ResponseURL,
			Text:       fmt.Sprintf("something wrong: %v\nplease contact the developer.", err),
		})
		return "", nil
	}
	return "", nil
}

func (h *SlashCommandHandler) handle(ctx context.Context, slash *slack.SlashCommand) error {
	text := strings.TrimSpace(slash.Text)
	if text == "" || text == "help" {
		return h.handleHelp(ctx, slash)
	}
	if text == "list" {
		return h.handleList(ctx, slash)
	}
	return nil
}

func (h *SlashCommandHandler) handleHelp(ctx context.Context, slash *slack.SlashCommand) error {
	h.cfg.PostSlackWebhook(ctx, &service.PostSlackWebhookInput{
		WebhookURL: slash.ResponseURL,
		Text:       "TODO: help message",
	})
	return nil
}

func (h *SlashCommandHandler) handleList(ctx context.Context, slash *slack.SlashCommand) error {
	permissions, err := h.cfg.GetSlackPermission(ctx, &repository.GetSlackPermissionInput{
		TeamID:    slash.TeamID,
		ChannelID: slash.ChannelID,
	})
	if err != nil {
		return err
	}

	var text strings.Builder
	for _, permission := range permissions.Repos {
		text.WriteString(permission)
		text.WriteString("\n")
	}
	h.cfg.PostSlackWebhook(ctx, &service.PostSlackWebhookInput{
		WebhookURL: slash.ResponseURL,
		Text:       text.String(),
	})
	return nil
}
