// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package ai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func TestGenerateSendsOpenAICompatibleRequest(ts *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			ts.Errorf("expected POST, got %s", request.Method)
		}
		if request.URL.Path != "/v1/chat/completions" {
			ts.Errorf("unexpected path: %s", request.URL.Path)
		}
		if got := request.Header.Get("Authorization"); got != "Bearer test-key" {
			ts.Errorf("unexpected authorization: %q", got)
		}

		var body chatCompletionRequest
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			ts.Errorf("decode request: %v", err)
		}
		if body.Model != "local-model" || len(body.Messages) != 2 {
			ts.Errorf("unexpected request: %#v", body)
		}
		if body.Messages[1].Content != "make a sequence" {
			ts.Errorf("unexpected user prompt: %q", body.Messages[1].Content)
		}
		if body.Messages[0].Content != systemPromptForProvider(t.AIProviderDefault, "make a sequence") {
			ts.Error("unexpected default system prompt")
		}
		if body.Temperature != 1 {
			ts.Errorf("expected temperature 1, got %v", body.Temperature)
		}
		if body.Stream != nil {
			ts.Errorf("expected stream to be omitted, got %v", *body.Stream)
		}

		_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"focus\n  tone 220 amplitude 10\n00:00:00 focus\n00:01:00 focus"}}]}`))
	}))
	defer server.Close()

	client, err := New(Config{APIKey: "test-key", BaseURL: server.URL, Model: "local-model", Temperature: 1})
	if err != nil {
		ts.Fatalf("New error: %v", err)
	}

	content, err := client.Generate(context.Background(), "make a sequence")
	if err != nil {
		ts.Fatalf("Generate error: %v", err)
	}
	if !strings.HasSuffix(content, "\n") || !strings.Contains(content, "tone 220") {
		ts.Fatalf("unexpected content: %q", content)
	}
}

func TestAppleFoundationProviderRequestsNonStreamingResponses(ts *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		if got := request.Header.Get("Authorization"); got != "Bearer test-key" {
			ts.Errorf("unexpected authorization: %q", got)
		}

		var body chatCompletionRequest
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			ts.Errorf("decode request: %v", err)
		}
		if body.Model != "configured-model" || body.Stream == nil || *body.Stream {
			ts.Errorf("unexpected Apple Foundation request: %#v", body)
		}
		if body.Messages[0].Content != appleFoundationSystemPrompt {
			ts.Error("unexpected Apple Foundation system prompt")
		}

		_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"focus\n  tone 220 amplitude 10\n00:00:00 focus\n00:01:00 focus"}}]}`))
	}))
	defer server.Close()

	client, err := New(Config{
		APIKey:      "test-key",
		BaseURL:     server.URL,
		Model:       "configured-model",
		Provider:    t.AIProviderAppleFoundation,
		Temperature: 0.2,
	})
	if err != nil {
		ts.Fatalf("New error: %v", err)
	}
	if _, err := client.Generate(context.Background(), "make a sequence"); err != nil {
		ts.Fatalf("Generate error: %v", err)
	}
	if _, err := client.Repair(context.Background(), "make a sequence", "invalid", errors.New("invalid SPSQ")); err != nil {
		ts.Fatalf("Repair error: %v", err)
	}
	if requests != 2 {
		ts.Fatalf("expected generation and repair requests, got %d", requests)
	}
}

func TestGenerateSendsConfiguredTemperature(ts *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var body chatCompletionRequest
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			ts.Errorf("decode request: %v", err)
		}
		if body.Temperature != 0.7 {
			ts.Errorf("unexpected temperature: %#v", body.Temperature)
		}

		_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"focus\n  tone 220 amplitude 10\n00:00:00 focus\n00:01:00 focus"}}]}`))
	}))
	defer server.Close()

	temperature := 0.7
	client, err := New(Config{
		APIKey:      "test-key",
		BaseURL:     server.URL,
		Model:       "local-model",
		Temperature: temperature,
	})
	if err != nil {
		ts.Fatalf("New error: %v", err)
	}
	if _, err := client.Generate(context.Background(), "make a sequence"); err != nil {
		ts.Fatalf("Generate error: %v", err)
	}
}

func TestGenerateReportsAPIError(ts *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte(`{"error":{"message":"invalid key"}}`))
	}))
	defer server.Close()

	client, err := New(Config{APIKey: "test-key", BaseURL: server.URL, Model: "local-model"})
	if err != nil {
		ts.Fatalf("New error: %v", err)
	}

	_, err = client.Generate(context.Background(), "make a sequence")
	if err == nil || !strings.Contains(err.Error(), "HTTP 401: invalid key") {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestNewRequiresAPIKey(ts *testing.T) {
	_, err := New(Config{Model: "test"})
	if err == nil || !strings.Contains(err.Error(), "SYNAPSEQ_AI_API_KEY") {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestGenerateHonorsContextCancellation(ts *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(50 * time.Millisecond)
	}))
	defer server.Close()

	client, err := New(Config{APIKey: "test-key", BaseURL: server.URL, Model: "local-model"})
	if err != nil {
		ts.Fatalf("New error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = client.Generate(ctx, "make a sequence")
	if !errors.Is(err, context.DeadlineExceeded) {
		ts.Fatalf("expected context deadline error, got %v", err)
	}
}
