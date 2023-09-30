package handler

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
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
	repository.SlackPermissionAllower
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
	args := regexp.MustCompile(`\s+`).Split(text, -1)
	if len(args) == 0 {
		return h.handleHelp(ctx, slash)
	}
	switch args[0] {
	case "help":
		return h.handleHelp(ctx, slash)
	case "list":
		return h.handleList(ctx, slash)
	case "allow":
		return h.handleAllow(ctx, slash, args[1:])
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

	if len(permissions.Repos) == 0 {
		h.cfg.PostSlackWebhook(ctx, &service.PostSlackWebhookInput{
			WebhookURL: slash.ResponseURL,
			Text:       "no repositories",
		})
		return nil
	}

	var text strings.Builder
	for _, permission := range permissions.Repos {
		text.WriteString(permission)
		text.WriteString("\n")
	}
	_, err = h.cfg.PostSlackWebhook(ctx, &service.PostSlackWebhookInput{
		WebhookURL: slash.ResponseURL,
		Text:       text.String(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *SlashCommandHandler) handleAllow(ctx context.Context, slash *slack.SlashCommand, args []string) error {
	_, err := h.cfg.AllowSlackPermission(ctx, &repository.AllowSlackPermissionInput{
		TeamID:    slash.TeamID,
		ChannelID: slash.ChannelID,
		Repos:     args,
	})
	if err != nil {
		return err
	}
	return nil
}
