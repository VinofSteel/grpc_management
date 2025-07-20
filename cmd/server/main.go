package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/vinofsteel/grpc-management/internal/handlers/user"
	"github.com/vinofsteel/grpc-management/pkg"
	"google.golang.org/grpc"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := fmt.Sprintf(":%s", port)

	slog.LogAttrs(
		ctx,
		slog.LevelInfo,
		"Starting TCP server",
	)

	lis, err := net.Listen("tcp", addr)
	if err != nil && err != http.ErrServerClosed {
		slog.ErrorContext(ctx, "TCP Server error", "error", err)
		os.Exit(1)
	}

	slog.LogAttrs(
		ctx,
		slog.LevelInfo,
		"Starting gRPC server",
		slog.Group("server", slog.String("address", addr)),
	)

	s := user.Server{}
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		slog.ErrorContext(ctx, "gRPC Server error", "error", err)
		os.Exit(1)
	}
}
