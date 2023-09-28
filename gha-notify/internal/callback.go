package internal

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/shogo82148/aws-xray-yasdk-go/xray"
	"github.com/slack-go/slack"
)

type CallbackHandler struct {
	appID        string
	clientID     string
	clientSecret string
}

func NewCallbackHandler(ctx context.Context) (*CallbackHandler, error) {
	ctx, seg := xray.BeginDummySegment(ctx)
	defer seg.Close()

	// initialize the AWS SDK
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	svc := ssm.NewFromConfig(cfg)

	// load the parameters
	appIDParam, err := svc.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String("/slack/app_id"),
	})
	if err != nil {
		return nil, err
	}
	clientIDParam, err := svc.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String("/slack/client_id"),
	})
	if err != nil {
		return nil, err
	}
	clientSecretParam, err := svc.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String("/slack/client_secret"),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	return &CallbackHandler{
		appID:        aws.ToString(appIDParam.Parameter.Value),
		clientID:     aws.ToString(clientIDParam.Parameter.Value),
		clientSecret: aws.ToString(clientSecretParam.Parameter.Value),
	}, nil
}

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: validate the request

	code := r.URL.Query().Get("code")
	resp, err := slack.GetOAuthV2ResponseContext(ctx, http.DefaultClient, h.clientID, h.clientSecret, code, "")
	if err != nil {
		slog.ErrorContext(ctx, "failed to get OAuth response", slog.String("error", err.Error()))
	}

	data, err := json.Marshal(resp)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get OAuth response", slog.String("error", err.Error()))
	}
	slog.InfoContext(ctx, "got OAuth response", slog.String("response", string(data)))
}
