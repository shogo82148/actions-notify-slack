package internal

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/shogo82148/aws-xray-yasdk-go/xray"
	"github.com/slack-go/slack"
)

const SlackAccessTokenTable = "slack-access-token"

type CallbackHandler struct {
	appID        string
	clientID     string
	clientSecret string

	dynamodb *dynamodb.Client
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
		dynamodb:     dynamodb.NewFromConfig(cfg),
	}, nil
}

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	now := time.Now()

	// TODO: validate the request

	code := r.URL.Query().Get("code")
	resp, err := slack.GetOAuthV2ResponseContext(ctx, http.DefaultClient, h.clientID, h.clientSecret, code, "")
	if err != nil {
		slog.ErrorContext(ctx, "failed to get OAuth response", slog.String("error", err.Error()))
	}

	expresAt := now.Add(time.Duration(resp.ExpiresIn) * time.Second)
	h.dynamodb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(SlackAccessTokenTable),
		Item: map[string]types.AttributeValue{
			"team_id": &types.AttributeValueMemberS{
				Value: resp.Team.ID,
			},
			"bot_user_id": &types.AttributeValueMemberS{
				Value: resp.BotUserID,
			},
			"access_token": &types.AttributeValueMemberS{
				Value: resp.AccessToken,
			},
			"scope": &types.AttributeValueMemberS{
				Value: resp.Scope,
			},
			"refresh_token": &types.AttributeValueMemberS{
				Value: resp.RefreshToken,
			},
			"expires_at": &types.AttributeValueMemberN{
				Value: strconv.FormatInt(expresAt.Unix(), 10),
			},
		},
	})
}
