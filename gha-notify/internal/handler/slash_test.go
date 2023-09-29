package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository"
	database "github.com/shogo82148/actions-notify-slack/gha-notify/internal/repository/mock"
	service "github.com/shogo82148/actions-notify-slack/gha-notify/internal/service/mock"
	"github.com/slack-go/slack"
)

func TestSlashHandler(t *testing.T) {
	const dummySigningSecret = "dummy-secret"
	var invoked bool
	invoker := service.LambdaInvokerFunc(func(ctx context.Context, input *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
		var s slack.SlashCommand
		if err := json.Unmarshal(input.Payload, &s); err != nil {
			return nil, err
		}
		if s.Command != "/notify-slack" {
			t.Errorf("unexpected command: want %q, got %q", "/notify-slack", s.Command)
		}
		if s.Text != "hello" {
			t.Errorf("unexpected text: want %q, got %q", "hello", s.Text)
		}
		if s.ResponseURL != "https://example.com" {
			t.Errorf("unexpected response_url: want %q, got %q", "https://example.com", s.ResponseURL)
		}
		invoked = true
		return &lambda.InvokeOutput{
			StatusCode: 200,
		}, nil
	})
	getter := database.SlackSigningSecretGetterFunc(func(ctx context.Context, input *repository.GetSlackSigningSecretInput) (*repository.GetSlackSigningSecretOutput, error) {
		return &repository.GetSlackSigningSecretOutput{
			SlackSigningSecret: dummySigningSecret,
		}, nil
	})

	h, err := NewSlashHandler(&SlashHandlerConfig{
		LambdaInvoker:            invoker,
		SlackSigningSecretGetter: getter,
		SlashFunctionName:        "slash-function",
	})
	if err != nil {
		t.Fatal(err)
	}

	// calculate the signature
	now := time.Now()
	body := `command=/notify-slack&text=hello&response_url=https://example.com`
	hash := hmac.New(sha256.New, []byte(dummySigningSecret))
	hash.Write([]byte("v0:" + strconv.FormatInt(now.Unix(), 10) + ":" + body))

	// build the request
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("X-Slack-Signature", "v0="+hex.EncodeToString(hash.Sum(nil)))
	req.Header.Set("X-Slack-Request-Timestamp", strconv.FormatInt(now.Unix(), 10))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(rw, req)

	// check the response
	resp := rw.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if !invoked {
		t.Error("lambda function is not invoked")
	}
}
