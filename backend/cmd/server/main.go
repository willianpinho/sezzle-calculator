package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/willianpinho/sezzle-calculator/backend/internal/httpapi"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	port := env("PORT", "8080")
	staticDir := env("STATIC_DIR", "../frontend/dist")
	if err := run(ctx, ":"+port, staticDir); err != nil {
		log.Fatal("server failed: ", err)
	}
}

func run(ctx context.Context, address, staticDir string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return serve(ctx, newServer(staticDir), listener)
}

func newServer(staticDir string) *http.Server {
	static := http.FileServer(http.Dir(filepath.Clean(staticDir)))
	return &http.Server{Handler: httpapi.New(static), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
}

func serve(ctx context.Context, server *http.Server, listener net.Listener) error {
	return serveWithTimeout(ctx, server, listener, 10*time.Second)
}

func serveWithTimeout(ctx context.Context, server *http.Server, listener net.Listener, timeout time.Duration) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Serve(listener)
	}()
	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	<-errCh // net/http guarantees ErrServerClosed after a successful Shutdown.
	return nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
