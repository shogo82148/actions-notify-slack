package main

import (
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/handler"
	"github.com/shogo82148/aws-xray-yasdk-go/xray/xraylog"
	"github.com/shogo82148/aws-xray-yasdk-go/xray/xrayslog"
)

var logger *slog.Logger

func init() {
	// initialize the logger
	h1 := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	h2 := xrayslog.NewHandler(h1, "trace_id")
	logger = slog.New(h2)
	slog.SetDefault(logger)
	xraylog.SetLogger(xrayslog.NewXRayLogger(h2))
}

func main() {
	h, err := handler.NewSlashCommandHandler(&handler.SlashCommandHandlerConfig{})
	if err != nil {
		slog.Error("failed to initialize the handler", slog.String("error", err.Error()))
		os.Exit(1)
	}
	lambda.Start(h.Handle)
}
