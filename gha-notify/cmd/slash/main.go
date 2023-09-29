package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/external"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/handler"
	"github.com/shogo82148/aws-xray-yasdk-go/xray/xraylog"
	"github.com/shogo82148/aws-xray-yasdk-go/xray/xrayslog"
	"github.com/shogo82148/aws-xray-yasdk-go/xrayhttp"
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
	h, err := newHandler(context.Background())
	if err != nil {
		slog.Error("failed to initialize the handler", slog.String("error", err.Error()))
		os.Exit(1)
	}
	lambda.Start(h.Handle)
}

func newHandler(ctx context.Context) (*handler.SlashCommandHandler, error) {
	httpClient := xrayhttp.Client(nil)
	svcSlack, err := external.NewSlack(&external.SlackConfig{
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, err
	}

	return handler.NewSlashCommandHandler(&handler.SlashCommandHandlerConfig{
		SlackWebhookPoster: svcSlack,
	})
}
