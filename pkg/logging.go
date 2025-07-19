package pkg

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func SetupSlog(ctx context.Context) {
	env := os.Getenv("ENV")

	logsDir := "logs"
	// #nosec G301
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		fmt.Printf("Error creating logs directory: %v\n", err)
		os.Exit(1)
	}

	// Create log file with timestamp in filename
	timestamp := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(logsDir, fmt.Sprintf("%s.log", timestamp))

	// #nosec G304 G302
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}

	// Create a multi-writer that writes to both stdout and the log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	var handler slog.Handler

	switch env {
	case "production":
		// Use JSON format in production for better parsing by log aggregation tools
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			// Add timestamp with consistent ISO8601 format
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.String(slog.TimeKey, a.Value.Time().Format(time.RFC3339))
				}
				return a
			},
		})
	case "test":
		// Minimal logging in test environment
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})
	default:
		// Development environment with more verbose logging
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true, // Include source file and line in logs for development
		})
	}

	// Replace the default slog logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Log the initial configuration
	logger.InfoContext(ctx, "Logger initialized",
		"environment", env,
		"level", logger.Handler().Enabled(ctx, slog.LevelInfo),
		"logFile", logFilePath,
	)
}

func ScheduleLogRotation(ctx context.Context) {
	// Get retention period from env or use default
	retentionDays := 7
	if retentionStr := os.Getenv("LOG_RETENTION_DAYS"); retentionStr != "" {
		if days, err := strconv.Atoi(retentionStr); err == nil && days > 0 {
			retentionDays = days
		}
	}

	slog.InfoContext(ctx, "Setting up log rotation", "retentionDays", retentionDays)

	// Immediately run log rotation on startup
	rotateLogFiles(ctx, retentionDays)

	// Calculate time until next 1:00 AM
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, now.Location())
	if now.After(nextRun) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	for {
		timer := time.NewTimer(time.Until(nextRun))
		select {
		case <-timer.C:
			rotateLogFiles(ctx, retentionDays)
			nextRun = nextRun.Add(24 * time.Hour)
		case <-ctx.Done():
			timer.Stop()
			return
		}
	}
}

// Utilities
func rotateLogFiles(ctx context.Context, retentionDays int) {
	logsDir := "logs"

	// Ensure logs directory exists
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		return
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	entries, err := os.ReadDir(logsDir)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read logs directory for rotation", "error", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".log" {
			// Extract date from filename (format: 2006-01-02.log)
			fileDate := strings.TrimSuffix(entry.Name(), ".log")
			parsedTime, err := time.Parse("2006-01-02", fileDate)

			if err != nil {
				slog.WarnContext(ctx, "Could not parse log filename for rotation",
					"filename", entry.Name(), "error", err)
				continue
			}

			// If file is older than retention period, delete it
			if parsedTime.Before(cutoffTime) {
				filePath := filepath.Join(logsDir, entry.Name())
				if err := os.Remove(filePath); err != nil {
					slog.ErrorContext(ctx, "Failed to delete old log file",
						"file", filePath, "error", err)
				} else {
					slog.InfoContext(ctx, "Deleted old log file", "file", filePath)
				}
			}
		}
	}
}
