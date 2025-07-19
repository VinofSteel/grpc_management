package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/vinofsteel/grpc-management/pkg"
)

func main() {
	// Creating application context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Setting up slog
	pkg.SetupSlog(ctx)
	daysToKeepLogs := 30
	go pkg.ScheduleLogRotation(ctx, daysToKeepLogs)

	// Initializing environment variables
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			slog.ErrorContext(ctx, "Error loading .env file", "error", err)
			os.Exit(1)
		}
	}

	// Handle common termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		slog.InfoContext(ctx, "Received cancellation signal", "signal", sig)
		cancel()
	}()
}
