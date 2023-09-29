package main

import (
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/handler"
)

func main() {
	h, err := handler.NewSlashCommandHandler(&handler.SlashCommandHandlerConfig{})
	if err != nil {
		slog.Error("failed to initialize the handler", slog.String("error", err.Error()))
		os.Exit(1)
	}
	lambda.Start(h.Handle)
}
