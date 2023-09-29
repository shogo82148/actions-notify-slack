package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/shogo82148/actions-notify-slack/gha-notify/internal"
	"github.com/shogo82148/aws-xray-yasdk-go/xray/xrayslog"
	httplogger "github.com/shogo82148/go-http-logger"
	"github.com/shogo82148/ridgenative"
)

var logger *slog.Logger

func init() {
	// initialize the logger
	h1 := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	h2 := xrayslog.NewHandler(h1, "trace_id")
	logger = slog.New(h2)
	slog.SetDefault(logger)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, World!\n")
	})

	callback, err := internal.NewCallbackHandler(context.Background())
	if err != nil {
		slog.Error("failed to initialize the callback handler", slog.String("error", err.Error()))
		os.Exit(1)
	}
	mux.Handle("/callback", callback)

	notifyHandler, err := internal.NewNotifyHandler(context.Background())
	if err != nil {
		slog.Error("failed to initialize the notify handler", slog.String("error", err.Error()))
		os.Exit(1)
	}
	mux.Handle("/notify", notifyHandler)

	logger := httplogger.NewSlogLogger(slog.LevelInfo, "http access log", logger)
	err = ridgenative.ListenAndServe(":8080", httplogger.LoggingHandler(logger, mux))
	if err != nil {
		slog.Error("failed to listen and serve: %v", err)
		os.Exit(1)
	}
}
