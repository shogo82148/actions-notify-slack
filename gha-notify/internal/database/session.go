package database

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
)

var _ repository.SessionPutter = (*SessionTable)(nil)
var _ repository.SessionGetter = (*SessionTable)(nil)

type SessionTable struct {
	cfg *SessionTableConfig
}

type SessionTableConfig struct {
	service.DynamoDBItemPutter
	service.DynamoDBItemGetter
	service.DynamoDBItemDeleter
	TableName string
}

func NewSessionTable(cfg *SessionTableConfig) (*SessionTable, error) {
	return &SessionTable{
		cfg: cfg,
	}, nil
}

// PutSession puts a session.
func (t *SessionTable) PutSession(ctx context.Context, input *repository.PutSessionInput) (*repository.PutSessionOutput, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	item := map[string]types.AttributeValue{
		"session_id": &types.AttributeValueMemberS{
			Value: input.SessionID,
		},
		"expires_at": &types.AttributeValueMemberN{
			Value: strconv.FormatFloat(timeToUnixTime(expiresAt), 'f', -1, 64),
		},
	}
	if input.State != "" {
		item["state"] = &types.AttributeValueMemberS{
			Value: input.State,
		}
	}
	if input.TeamID != "" {
		item["team_id"] = &types.AttributeValueMemberS{
			Value: input.TeamID,
		}
	}
	if input.TeamName != "" {
		item["team_name"] = &types.AttributeValueMemberS{
			Value: input.TeamName,
		}
	}

	_, err := t.cfg.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(t.cfg.TableName),
		Item:      item,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// GetSession gets a session.
func (t *SessionTable) GetSession(ctx context.Context, input *repository.GetSessionInput) (*repository.GetSessionOutput, error) {
	out, err := t.cfg.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(t.cfg.TableName),
		Key: map[string]types.AttributeValue{
			"session_id": &types.AttributeValueMemberS{
				Value: input.SessionID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if _, ok := out.Item["session_id"]; !ok {
		return &repository.GetSessionOutput{}, nil
	}
	conv := new(attrConverter)
	sessionID := conv.convertString(out.Item["session_id"])

	var state string
	if attr, ok := out.Item["state"]; ok {
		state = conv.convertString(attr)
	}
	var teamID string
	if attr, ok := out.Item["team_id"]; ok {
		teamID = conv.convertString(attr)
	}
	var teamName string
	if attr, ok := out.Item["team_name"]; ok {
		teamName = conv.convertString(attr)
	}
	if conv.err != nil {
		return nil, conv.err
	}
	return &repository.GetSessionOutput{
		SessionID: sessionID,
		State:     state,
		TeamID:    teamID,
		TeamName:  teamName,
	}, nil
}

// DeleteSession deletes a session.
func (t *SessionTable) DeleteSession(ctx context.Context, input *repository.DeleteSessionInput) (*repository.DeleteSessionOutput, error) {
	_, err := t.cfg.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(t.cfg.TableName),
		Key: map[string]types.AttributeValue{
			"session_id": &types.AttributeValueMemberS{
				Value: input.SessionID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
