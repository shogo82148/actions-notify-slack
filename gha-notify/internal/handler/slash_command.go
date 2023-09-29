package handler

import (
	"context"

	"github.com/slack-go/slack"
)

type SlashCommandHandler struct {
	cfg *SlashCommandHandlerConfig
}

type SlashCommandHandlerConfig struct {
}

func NewSlashCommandHandler(cfg *SlashCommandHandlerConfig) (*SlashCommandHandler, error) {
	return &SlashCommandHandler{
		cfg: cfg,
	}, nil
}

func (h *SlashCommandHandler) Handle(ctx context.Context, slash *slack.SlashCommand) (string, error) {
	return "", nil
}
