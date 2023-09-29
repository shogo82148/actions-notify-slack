package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/service"
)

var _ service.LambdaInvoker = LambdaInvokerFunc(nil)

type LambdaInvokerFunc func(ctx context.Context, input *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)

func (f LambdaInvokerFunc) Invoke(ctx context.Context, input *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
	return f(ctx, input, optFns...)
}
