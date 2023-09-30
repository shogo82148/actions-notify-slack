package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
)

var _ repository.SlackPermissionGetter = (*SlackPermissionTable)(nil)

type SlackPermissionTable struct {
	cfg *SlackPermissionTableConfig
}

type SlackPermissionTableConfig struct {
	service.DynamoDBItemPutter
	service.DynamoDBItemGetter
	service.DynamoDBItemUpdater
	TableName string
}

func NewSlackPermissionTable(cfg *SlackPermissionTableConfig) (*SlackPermissionTable, error) {
	return &SlackPermissionTable{cfg: cfg}, nil
}

func (t *SlackPermissionTable) GetSlackPermission(ctx context.Context, input *repository.GetSlackPermissionInput) (*repository.GetSlackPermissionOutput, error) {
	out, err := t.cfg.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(t.cfg.TableName),
		Key: map[string]types.AttributeValue{
			"team_id": &types.AttributeValueMemberS{
				Value: input.TeamID,
			},
			"channel_id": &types.AttributeValueMemberS{
				Value: input.ChannelID,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if _, ok := out.Item["repos"]; !ok {
		return &repository.GetSlackPermissionOutput{
			TeamID:    input.TeamID,
			ChannelID: input.ChannelID,
			Repos:     []string{},
		}, nil
	}

	conv := new(attrConverter)
	teamID := conv.convertString(out.Item["team_id"])
	channelID := conv.convertString(out.Item["channel_id"])
	repos := conv.convertStringSet(out.Item["repos"])
	if conv.err != nil {
		return nil, conv.err
	}

	return &repository.GetSlackPermissionOutput{
		TeamID:    teamID,
		ChannelID: channelID,
		Repos:     repos,
	}, nil
}

func (t *SlackPermissionTable) AllowSlackPermission(ctx context.Context, input *repository.AllowSlackPermissionInput) (*repository.AllowSlackPermissionOutput, error) {
	_, err := t.cfg.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(t.cfg.TableName),
		Key: map[string]types.AttributeValue{
			"team_id": &types.AttributeValueMemberS{
				Value: input.TeamID,
			},
			"channel_id": &types.AttributeValueMemberS{
				Value: input.ChannelID,
			},
		},
		UpdateExpression: aws.String("ADD repos :repo"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":repo": &types.AttributeValueMemberSS{
				Value: input.Repos,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &repository.AllowSlackPermissionOutput{}, nil
}
