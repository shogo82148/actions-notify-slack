.PHONY: help
help:
	cat Makefile

.PHONY: test
test:
	go test -v ./...

.PHONY: build-NotifyFunction
build-NotifyFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(ARTIFACTS_DIR)/bootstrap -tags lambda.norpc ./cmd/notify

.PHONY: build-SlashFunction
build-SlashFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(ARTIFACTS_DIR)/bootstrap -tags lambda.norpc ./cmd/slash
