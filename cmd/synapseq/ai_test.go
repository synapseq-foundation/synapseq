// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

func TestRunAIWritesValidatedSequence(ts *testing.T) {
	server := aiTestServer(ts, `relax
  tone 220 amplitude 10
00:00:00 relax
00:10:00 relax`)
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	outputPath := filepath.Join(ts.TempDir(), "relax.spsq")
	var status bytes.Buffer
	err := runAI(
		context.Background(),
		"Generate a 10 minute session",
		[]string{outputPath},
		&cli.CLIOptions{AIBaseURL: server.URL},
		&status,
		&bytes.Buffer{},
	)
	if err != nil {
		ts.Fatalf("runAI error: %v", err)
	}
	content, err := os.ReadFile(outputPath)
	if err != nil {
		ts.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(content), "tone 220") {
		ts.Fatalf("unexpected output: %q", content)
	}
	if !strings.Contains(status.String(), "Generated:") {
		ts.Fatalf("unexpected status: %q", status.String())
	}
	if !strings.Contains(status.String(), aiProgressMessage) {
		ts.Fatalf("expected progress message, got %q", status.String())
	}
	generatedIndex := strings.Index(status.String(), "Generated:")
	if generatedIndex < 0 || !strings.Contains(status.String()[:generatedIndex], "\n") {
		ts.Fatalf("expected generated status on a new line, got %q", status.String())
	}
}

func TestRunAIStreamsOnlySPSQ(ts *testing.T) {
	server := aiTestServer(ts, `focus
  noise pink smooth 20 amplitude 10
00:00:00 focus
00:05:00 focus`)
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	var output bytes.Buffer
	err := runAI(context.Background(), "Generate a session", []string{"-"}, &cli.CLIOptions{AIBaseURL: server.URL}, &bytes.Buffer{}, &output)
	if err != nil {
		ts.Fatalf("runAI error: %v", err)
	}
	if strings.Contains(output.String(), "Generated:") || !strings.Contains(output.String(), "noise pink") {
		ts.Fatalf("unexpected standard output: %q", output.String())
	}
}

func TestRunAIRejectsExistingOutputBeforeRequest(ts *testing.T) {
	outputPath := filepath.Join(ts.TempDir(), "existing.spsq")
	if err := os.WriteFile(outputPath, []byte("existing"), 0o644); err != nil {
		ts.Fatalf("create output: %v", err)
	}

	err := runAI(context.Background(), "Generate focus", []string{outputPath}, &cli.CLIOptions{}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAIRejectsEmptyPrompt(ts *testing.T) {
	err := runAI(context.Background(), "", []string{"-"}, &cli.CLIOptions{}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "prompt cannot be empty") {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestAIOutputPathUsesIntentAndDuration(ts *testing.T) {
	if got := aiOutputPath("Generate a 10 minutes of relaxation sequence"); got != "relaxation-10m.spsq" {
		ts.Fatalf("unexpected output path: %q", got)
	}
}

func TestCLIAITemperature(ts *testing.T) {
	temperature, err := cliAITemperature("0.2")
	if err != nil {
		ts.Fatalf("cliAITemperature error: %v", err)
	}
	if temperature == nil || *temperature != 0.2 {
		ts.Fatalf("unexpected temperature: %#v", temperature)
	}

	if _, err := cliAITemperature("not-a-number"); err == nil {
		ts.Fatal("expected invalid temperature error")
	}
}

func TestCLIAITimeout(ts *testing.T) {
	timeout, err := cliAITimeout("90s")
	if err != nil {
		ts.Fatalf("cliAITimeout error: %v", err)
	}
	if timeout == nil || *timeout != 90*time.Second {
		ts.Fatalf("unexpected timeout: %#v", timeout)
	}

	if _, err := cliAITimeout("not-a-duration"); err == nil {
		ts.Fatal("expected invalid timeout error")
	}
}

func TestStartAIProgressDoesNotAnimateNonTerminal(ts *testing.T) {
	var output bytes.Buffer
	progress := startAIProgress(&output, false)
	if progress != nil {
		ts.Fatal("expected no spinner for a non-terminal writer")
	}
	if got := output.String(); !strings.Contains(got, aiProgressMessage) {
		ts.Fatalf("expected progress message, got %q", got)
	}
}

func TestStartAIProgressRespectsQuiet(ts *testing.T) {
	var output bytes.Buffer
	progress := startAIProgress(&output, true)
	if progress != nil {
		ts.Fatal("expected no spinner in quiet mode")
	}
	if got := output.String(); got != "" {
		ts.Fatalf("expected no status output, got %q", got)
	}
}

func TestAIProgressStopEndsLine(ts *testing.T) {
	var output bytes.Buffer
	progress := &aiProgress{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go progress.animate(&output)
	progress.Stop()

	if !strings.HasSuffix(output.String(), "\n") {
		ts.Fatalf("expected progress output to end with a newline, got %q", output.String())
	}
}

func aiTestServer(ts *testing.T, content string) *httptest.Server {
	ts.Helper()
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":` + strconv.Quote(content) + `}}]}`))
	}))
}
