package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
)

func hello(ctx context.Context, slash *slack.SlashCommand) (string, error) {
	slog.InfoContext(ctx, "slash command", slog.String("command", slash.Command), slog.String("text", slash.Text))
	return "Hello λ!", nil
}

func main() {
	lambda.Start(hello)
}
