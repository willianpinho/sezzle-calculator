package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestServeGracefulShutdown(t *testing.T) {
	staticDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(staticDir, "index.html"), []byte("calculator"), 0o600); err != nil {
		t.Fatal(err)
	}
	listener := listen(t)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- serve(ctx, newServer(staticDir), listener) }()
	defer func() {
		cancel()
		select {
		case <-done:
		default:
		}
	}()
	response := waitForResponse(t, "http://"+listener.Addr().String()+"/api/health")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("health status = %d", response.StatusCode)
	}
	response.Body.Close()
	cancel()
	if err := <-done; err != nil {
		t.Fatalf("serve() = %v", err)
	}
}

func TestRunReportsListenFailure(t *testing.T) {
	if err := run(context.Background(), "not-an-address", t.TempDir()); err == nil {
		t.Fatal("run() error = nil; want listen error")
	}
}

func TestRunHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := run(ctx, "127.0.0.1:0", t.TempDir()); err != nil {
		t.Fatalf("run() = %v", err)
	}
}

func TestServeReportsListenerFailure(t *testing.T) {
	listener := listen(t)
	listener.Close()
	if err := serve(context.Background(), newServer(t.TempDir()), listener); err == nil {
		t.Fatal("serve() error = nil; want listener error")
	}
}

func TestServeHandlesServerClose(t *testing.T) {
	listener := listen(t)
	server := newServer(t.TempDir())
	defer server.Close()
	done := make(chan error, 1)
	go func() { done <- serve(context.Background(), server, listener) }()
	response := waitForResponse(t, "http://"+listener.Addr().String()+"/api/health")
	response.Body.Close()
	if err := server.Close(); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatalf("serve() = %v", err)
	}
}

func TestServeReportsShutdownTimeout(t *testing.T) {
	listener := listen(t)
	started, release := make(chan struct{}), make(chan struct{})
	server := &http.Server{Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) { close(started); <-release })}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		select {
		case <-release:
		default:
			close(release)
		}
	}()
	done := make(chan error, 1)
	go func() { done <- serveWithTimeout(ctx, server, listener, 10*time.Millisecond) }()
	requestDone := make(chan struct{})
	go func() {
		response, _ := (&http.Client{Timeout: time.Second}).Get("http://" + listener.Addr().String())
		if response != nil {
			response.Body.Close()
		}
		close(requestDone)
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("request did not reach handler")
	}
	cancel()
	var err error
	select {
	case err = <-done:
	case <-time.After(time.Second):
		t.Fatal("shutdown did not time out")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("serve() = %v; want deadline", err)
	}
	close(release)
	select {
	case <-requestDone:
	case <-time.After(time.Second):
		t.Fatal("request did not finish")
	}
}

func TestEnv(t *testing.T) {
	t.Setenv("CALCULATOR_TEST_VALUE", "configured")
	if got := env("CALCULATOR_TEST_VALUE", "fallback"); got != "configured" {
		t.Fatalf("env() = %q", got)
	}
	t.Setenv("CALCULATOR_TEST_UNSET", "")
	if got := env("CALCULATOR_TEST_UNSET", "fallback"); got != "fallback" {
		t.Fatalf("env() = %q", got)
	}
}

func TestMainHandlesSIGTERM(t *testing.T) {
	if os.Getenv("CALCULATOR_MAIN_SHUTDOWN_HELPER") == "1" {
		done := make(chan struct{})
		go func() { main(); close(done) }()
		response := waitForResponse(t, "http://127.0.0.1:"+os.Getenv("PORT")+"/api/health")
		response.Body.Close()
		if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
			t.Fatal(err)
		}
		select {
		case <-done:
			return
		case <-time.After(2 * time.Second):
			t.Fatal("main did not stop after SIGTERM")
		}
	}
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	_, port, err := net.SplitHostPort(listener.Addr().String())
	listener.Close()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestMainHandlesSIGTERM$")
	command.Env = testEnvironment("CALCULATOR_MAIN_SHUTDOWN_HELPER=1", "PORT="+port, "STATIC_DIR="+t.TempDir())
	if output, err := command.CombinedOutput(); ctx.Err() != nil || err != nil {
		t.Fatalf("shutdown helper failed: %v %s", err, output)
	}
}

func TestMainStartupFailureExits(t *testing.T) {
	if os.Getenv("CALCULATOR_MAIN_FAILURE_HELPER") == "1" {
		main()
		return
	}
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestMainStartupFailureExits$")
	command.Env = testEnvironment("CALCULATOR_MAIN_FAILURE_HELPER=1", "PORT="+port, "STATIC_DIR="+t.TempDir())
	output, err := command.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("startup helper timed out: %s", output)
	}
	var exitErr *exec.ExitError
	if err == nil || !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 || !strings.Contains(string(output), "server failed:") {
		t.Fatalf("startup helper = %v (%s); want exit 1 with startup failure", err, output)
	}
}

func listen(t *testing.T) net.Listener {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	return listener
}
func waitForResponse(t *testing.T, url string) *http.Response {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		response, err := http.Get(url)
		if err == nil {
			return response
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("server did not become ready")
	return nil
}

func testEnvironment(values ...string) []string {
	keys := map[string]bool{}
	for _, value := range values {
		keys[strings.SplitN(value, "=", 2)[0]] = true
	}
	result := make([]string, 0, len(os.Environ())+len(values))
	for _, value := range os.Environ() {
		if !keys[strings.SplitN(value, "=", 2)[0]] {
			result = append(result, value)
		}
	}
	return append(result, values...)
}
