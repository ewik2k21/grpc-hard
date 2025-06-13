package main

import (
	"github.com/ewik2k21/grpc-hard/cmd/server"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	server.Execute(logger)
}
