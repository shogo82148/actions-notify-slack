package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var _ SSMParameterGetter = (*ssm.Client)(nil)

// SSMParameterGetter is an interface for ssm.Client.GetParameter.
type SSMParameterGetter interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}
