package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type LambdaInvoker interface {
	Invoke(ctx context.Context, input *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}
