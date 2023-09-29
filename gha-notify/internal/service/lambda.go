package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// LambdaInvoker is an interface for lambda.Client.Invoke.
type LambdaInvoker interface {
	Invoke(ctx context.Context, input *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}
