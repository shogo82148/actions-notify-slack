package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var _ DynamoDBItemPutter = (*dynamodb.Client)(nil)

// DynamoDBItemPutter is an interface for dynamodb.Client.PutItem.
type DynamoDBItemPutter interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

var _ DynamoDBItemGetter = (*dynamodb.Client)(nil)

// DynamoDBItemGetter is an interface for dynamodb.Client.GetItem.
type DynamoDBItemGetter interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

var _ DynamoDBItemUpdater = (*dynamodb.Client)(nil)

// DynamoDBItemUpdater is an interface for dynamodb.Client.UpdateItem.
type DynamoDBItemUpdater interface {
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}
