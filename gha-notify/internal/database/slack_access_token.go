package database

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
)

var _ repository.SlackAccessTokenPutter = (*SlackAccessTokenTable)(nil)

type SlackAccessTokenTable struct {
	cfg *SlackAccessTokenTableConfig
}

type SlackAccessTokenTableConfig struct {
	service.DynamoDBItemPutter
	service.DynamoDBItemGetter
	TableName string
}

func NewSlackAccessTokenTable(cfg *SlackAccessTokenTableConfig) (*SlackAccessTokenTable, error) {
	return &SlackAccessTokenTable{
		cfg: cfg,
	}, nil
}

// PutSlackAccessToken puts a slack access token.
func (t *SlackAccessTokenTable) PutSlackAccessToken(ctx context.Context, input *repository.PutSlackAccessTokenInput) (*repository.PutSlackAccessTokenOutput, error) {
	_, err := t.cfg.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &t.cfg.TableName,
		Item: map[string]types.AttributeValue{
			"team_id": &types.AttributeValueMemberS{
				Value: input.TeamID,
			},
			"bot_user_id": &types.AttributeValueMemberS{
				Value: input.BotUserID,
			},
			"access_token": &types.AttributeValueMemberS{
				Value: input.AccessToken,
			},
			"scope": &types.AttributeValueMemberS{
				Value: input.Scope,
			},
			"refresh_token": &types.AttributeValueMemberS{
				Value: input.RefreshToken,
			},
			"expires_at": &types.AttributeValueMemberN{
				Value: strconv.FormatFloat(timeToUnixTime(input.ExpiresAt), 'f', -1, 64),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &repository.PutSlackAccessTokenOutput{}, nil
}

func (t *SlackAccessTokenTable) GetSlackAccessToken(ctx context.Context, input *repository.GetSlackAccessTokenInput) (*repository.GetSlackAccessTokenOutput, error) {
	out, err := t.cfg.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &t.cfg.TableName,
		Key: map[string]types.AttributeValue{
			"team_id": &types.AttributeValueMemberS{
				Value: input.TeamID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	conv := new(attrConverter)
	teamID := conv.convertString(out.Item["team_id"])
	botUserID := conv.convertString(out.Item["bot_user_id"])
	accessToken := conv.convertString(out.Item["access_token"])
	scope := conv.convertString(out.Item["scope"])
	refreshToken := conv.convertString(out.Item["refresh_token"])
	expiresAt := conv.convertNumber(out.Item["expires_at"])
	if conv.err != nil {
		return nil, conv.err
	}

	return &repository.GetSlackAccessTokenOutput{
		TeamID:       teamID,
		BotUserID:    botUserID,
		AccessToken:  accessToken,
		Scope:        scope,
		RefreshToken: refreshToken,
		ExpiresAt:    unixTimeToTime(expiresAt),
	}, nil
}
