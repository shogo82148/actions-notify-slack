.PHONY: help
help:
	cat Makefile

.PHONY: build-NotifyFunction
build-NotifyFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(ARTIFACTS_DIR)/notify -tags lambda.norpc ./cmd/notify
