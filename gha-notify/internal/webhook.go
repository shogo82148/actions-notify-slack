package internal

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/shogo82148/aws-xray-yasdk-go/xray"
	"github.com/slack-go/slack"
)

type Webhook struct {
	signingSecret string
	lambda        *lambda.Client
}

func NewWebhook(ctx context.Context) (*Webhook, error) {
	ctx, seg := xray.BeginDummySegment(ctx)
	defer seg.Close()

	// initialize the AWS SDK
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	svc := ssm.NewFromConfig(cfg)

	signingSecretParam, err := svc.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String("/slack/signing_secret"),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	return &Webhook{
		signingSecret: aws.ToString(signingSecretParam.Parameter.Value),
		lambda:        lambda.NewFromConfig(cfg),
	}, nil
}

func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	verifier, err := slack.NewSecretsVerifier(r.Header, h.signingSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cr := *r
	cr.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(&cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := verifier.Ensure(); err != nil {
		slog.ErrorContext(ctx, "failed to verify the request", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// handle the slash command
	payload, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	arn := os.Getenv("SLASH_FUNCTION_NAME")
	_, err = h.lambda.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(arn),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payload,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
