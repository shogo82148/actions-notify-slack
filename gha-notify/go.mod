module github.com/shogo82148/actions-notify-slack/gha-notify

go 1.21.1

require (
	github.com/aws/aws-lambda-go v1.41.0
	github.com/aws/aws-sdk-go-v2 v1.22.2
	github.com/aws/aws-sdk-go-v2/config v1.23.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.25.1
	github.com/aws/aws-sdk-go-v2/service/lambda v1.45.0
	github.com/aws/aws-sdk-go-v2/service/ssm v1.42.1
	github.com/shogo82148/aws-xray-yasdk-go v1.7.2
	github.com/shogo82148/aws-xray-yasdk-go/xrayaws-v2 v1.1.5
	github.com/shogo82148/go-http-logger v1.3.0
	github.com/shogo82148/goat v0.0.6
	github.com/shogo82148/memoize v0.0.4
	github.com/shogo82148/ridgenative v1.4.0
	github.com/slack-go/slack v0.12.3
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.5.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.8.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.17.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.19.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.25.1 // indirect
	github.com/aws/smithy-go v1.16.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/shogo82148/pointer v1.3.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
)
