package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/database"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/external"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/handler"
	"github.com/shogo82148/aws-xray-yasdk-go/xray/xraylog"
	"github.com/shogo82148/aws-xray-yasdk-go/xray/xrayslog"
	"github.com/shogo82148/aws-xray-yasdk-go/xrayaws-v2"
	"github.com/shogo82148/aws-xray-yasdk-go/xrayhttp"
	httplogger "github.com/shogo82148/go-http-logger"
	"github.com/shogo82148/ridgenative"
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
	ctx := context.Background()
	mux, err := NewMux(ctx)
	if err != nil {
		slog.Error("failed to initialize the mux", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger := httplogger.NewSlogLogger(slog.LevelInfo, "http access log", logger)
	err = ridgenative.ListenAndServe(":8080", httplogger.LoggingHandler(logger, mux))
	if err != nil {
		slog.Error("failed to listen and serve: %v", err)
		os.Exit(1)
	}
}

func NewMux(ctx context.Context) (http.Handler, error) {
	cfg, err := config.LoadDefaultConfig(ctx, xrayaws.WithXRay())
	if err != nil {
		return nil, err
	}
	svcDynamoDB := dynamodb.NewFromConfig(cfg)
	svcSSM := ssm.NewFromConfig(cfg)
	svcLambda := lambda.NewFromConfig(cfg)

	httpClient := xrayhttp.Client(nil)
	svcSlack, err := external.NewSlack(&external.SlackConfig{
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, err
	}

	svcGitHub, err := external.NewGitHub(&external.GitHubConfig{
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, err
	}

	params, err := database.NewParameters(&database.ParametersConfig{
		SSMParameterGetter: svcSSM,
	})
	if err != nil {
		return nil, err
	}

	slackAccessTokenTable, err := database.NewSlackAccessTokenTable(&database.SlackAccessTokenTableConfig{
		DynamoDBItemPutter: svcDynamoDB,
		DynamoDBItemGetter: svcDynamoDB,
		TableName:          "slack-access-token",
	})
	if err != nil {
		return nil, err
	}

	slackPermissionTable, err := database.NewSlackPermissionTable(&database.SlackPermissionTableConfig{
		DynamoDBItemPutter:  svcDynamoDB,
		DynamoDBItemGetter:  svcDynamoDB,
		DynamoDBItemUpdater: svcDynamoDB,
		TableName:           "slack-permission",
	})
	if err != nil {
		return nil, err
	}

	sessionTable, err := database.NewSessionTable(&database.SessionTableConfig{
		DynamoDBItemPutter:  svcDynamoDB,
		DynamoDBItemGetter:  svcDynamoDB,
		DynamoDBItemDeleter: svcDynamoDB,
		TableName:           "session",
	})
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	index, err := handler.NewIndexHandler(&handler.IndexHandlerConfig{
		SlackClientIDGetter: params,
		SessionGetter:       sessionTable,
		SessionPutter:       sessionTable,
	})
	if err != nil {
		return nil, err
	}
	mux.Handle("/", index)

	callback, err := handler.NewCallbackHandler(&handler.CallbackHandlerConfig{
		OAuthV2ResponseGetter:   svcSlack,
		SlackClientIDGetter:     params,
		SlackClientSecretGetter: params,
		SlackAccessTokenPutter:  slackAccessTokenTable,
		SessionGetter:           sessionTable,
		SessionPutter:           sessionTable,
		SessionDeleter:          sessionTable,
	})
	if err != nil {
		return nil, err
	}
	mux.Handle("/callback", callback)

	notifyHandler, err := handler.NewNotifyHandler(&handler.NotifyHandlerConfig{
		OAuthV2ResponseRefresher: svcSlack,
		SlackMessagePoster:       svcSlack,
		GitHubIDTokenParser:      svcGitHub,
		SlackClientIDGetter:      params,
		SlackClientSecretGetter:  params,
		SlackAccessTokenGetter:   slackAccessTokenTable,
		SlackAccessTokenPutter:   slackAccessTokenTable,
		SlackPermissionGetter:    slackPermissionTable,
	})
	if err != nil {
		return nil, err
	}
	mux.Handle("/notify", notifyHandler)

	slash, err := handler.NewSlashHandler(&handler.SlashHandlerConfig{
		SlackSigningSecretGetter: params,
		LambdaInvoker:            svcLambda,
		SlashFunctionName:        os.Getenv("SLASH_FUNCTION_NAME"),
	})
	if err != nil {
		return nil, err
	}
	mux.Handle("/slash", slash)

	return mux, nil
}
