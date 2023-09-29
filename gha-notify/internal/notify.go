package internal

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/slack-go/slack"
)

type NotifyHandler struct {
	dynamodb *dynamodb.Client
}

func NewNotifyHandler(ctx context.Context) (*NotifyHandler, error) {
	// initialize the AWS SDK
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &NotifyHandler{
		dynamodb: dynamodb.NewFromConfig(cfg),
	}, nil
}

func (h *NotifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: get the parameters from the request
	teamID := "T3G1HAY66"
	channelD := "C3GMGG162"

	// get the access token
	out, err := h.dynamodb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(SlackAccessTokenTable),
		Key: map[string]types.AttributeValue{
			"team_id": &types.AttributeValueMemberS{
				Value: teamID,
			},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: refresh token if needed
	token := out.Item["access_token"].(*types.AttributeValueMemberS).Value

	// send a message
	api := slack.New(token)
	_, _, err = api.PostMessageContext(ctx, channelD, slack.MsgOptionText("Hello, World!", false))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
