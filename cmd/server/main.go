package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/atran25/synckor/internal/api"
	"github.com/atran25/synckor/internal/config"
	"github.com/atran25/synckor/internal/database"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Load Config struct from .env file/environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Get a connection to the database
	databaseConnection, err := database.GetConnection()
	if err != nil {
		return err
	}

	h, err := api.NewServer(cfg, databaseConnection)
	if err != nil {
		return err
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h,
	}

	go func() {
		slog.Info("Server starting on", "Server Address", s.Addr)
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "Error", err)
		}
	}()

	// Graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := s.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server failed to shutdown", "Error", err)
		}
	}()
	wg.Wait()

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		slog.Error("Error running server", "Error", err)
		os.Exit(1)
	}
	slog.Info("Server stopped gracefully")
	os.Exit(0)
}
