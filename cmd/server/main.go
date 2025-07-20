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
	"time"

	"github.com/joho/godotenv"
	"github.com/vinofsteel/grpc-management/internal/handlers"
	"github.com/vinofsteel/grpc-management/internal/handlers/proto_user"
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

	handlers := handlers.Handlers{}
	grpcServer := grpc.NewServer()
	proto_user.RegisterUserServiceServer(grpcServer, &handlers)

	// Channel to capture server errors
	serverErrCh := make(chan error, 1)

	// Start gRPC server in a goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			serverErrCh <- err
		}
	}()

	// Handle common termination signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either a signal or server error
	select {
	case sig := <-sigCh:
		slog.InfoContext(ctx, "Received cancellation signal", "signal", sig)
	case err := <-serverErrCh:
		slog.ErrorContext(ctx, "gRPC Server error", "error", err)
		cancel()
		os.Exit(1)
	case <-ctx.Done():
		slog.InfoContext(ctx, "Context cancelled")
	}

	// Graceful shutdown
	slog.InfoContext(ctx, "Shutting down gRPC server...")

	// Create a context with timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Channel to signal when graceful stop is complete
	gracefulStopCh := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(gracefulStopCh)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-gracefulStopCh:
		slog.InfoContext(ctx, "gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		slog.WarnContext(ctx, "Graceful shutdown timeout, forcing stop")
		grpcServer.Stop() // Force stop
	}

	cancel()
	slog.InfoContext(ctx, "Application shutdown complete")
}
