// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAIValidatesGeneratedSPSQ(ts *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"relax\n  tone 220 amplitude 10\n00:00:00 relax\n00:10:00 relax"}}]}`))
	}))
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	loaded, err := NewAppContext().AI(context.Background(), "make sequence", testAIOptions(server.URL))
	if err != nil {
		ts.Fatalf("AI error: %v", err)
	}
	if got := string(loaded.RawContent()); !strings.Contains(got, "relax") {
		ts.Fatalf("unexpected content: %q", got)
	}
}

func TestAIRejectsInvalidSPSQ(ts *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"this is not spsq"}}]}`))
	}))
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	_, err := NewAppContext().AI(context.Background(), "make sequence", testAIOptions(server.URL))
	if err == nil || !strings.Contains(err.Error(), "AI did not understand the prompt") {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestAIOptionsOverrideEnvironment(ts *testing.T) {
	ts.Setenv("SYNAPSEQ_AI_MODEL", "environment-model")
	ts.Setenv("SYNAPSEQ_AI_BASE_URL", "https://environment.invalid")

	options := &AIOptions{Model: "flag-model", BaseURL: "https://flag.invalid"}
	if got := aiModel(options); got != "flag-model" {
		ts.Fatalf("expected flag model, got %q", got)
	}
	if got := aiBaseURL(options); got != "https://flag.invalid" {
		ts.Fatalf("expected flag base URL, got %q", got)
	}
}

func TestAIValidatesConfiguredTemperature(ts *testing.T) {
	if err := validateAITemperature(0.7); err != nil {
		ts.Fatalf("validateAITemperature error: %v", err)
	}
}

func TestAITemperatureRejectsInvalidValue(ts *testing.T) {
	if err := validateAITemperature(3); err == nil {
		ts.Fatal("expected invalid temperature error")
	}
}

func TestAIRepairsInvalidSPSQ(ts *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		if requests == 2 {
			var body struct {
				Messages []struct {
					Content string `json:"content"`
				} `json:"messages"`
			}
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				ts.Errorf("decode repair request: %v", err)
			}
			if len(body.Messages) != 2 || !strings.Contains(body.Messages[1].Content, "Validation error") {
				ts.Errorf("unexpected repair request: %#v", body.Messages)
			}
			_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"relax\n  tone 220 amplitude 10\n00:00:00 relax\n00:10:00 relax"}}]}`))
			return
		}

		_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"not valid spsq"}}]}`))
	}))
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	loaded, err := NewAppContext().AI(context.Background(), "make sequence", testAIOptions(server.URL))
	if err != nil {
		ts.Fatalf("AI error: %v", err)
	}
	if requests != 2 || !strings.Contains(string(loaded.RawContent()), "relax") {
		ts.Fatalf("unexpected repair result: requests=%d content=%q", requests, loaded.RawContent())
	}
}

func TestAIRepairsInvalidSemantics(ts *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		if requests == 2 {
			var body struct {
				Messages []struct {
					Content string `json:"content"`
				} `json:"messages"`
			}
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				ts.Errorf("decode repair request: %v", err)
			}
			if !strings.Contains(body.Messages[1].Content, "audible carrier between 100 and 600 Hz") {
				ts.Errorf("unexpected semantic repair request: %#v", body.Messages)
			}
			_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"gamma\n  tone 220 binaural 40 amplitude 10\n00:00:00 gamma\n00:05:00 gamma"}}]}`))
			return
		}

		_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"gamma\n  tone 40 binaural 40 amplitude 10\n00:00:00 gamma\n00:05:00 gamma"}}]}`))
	}))
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	loaded, err := NewAppContext().AI(context.Background(), "Create a gamma session", testAIOptions(server.URL))
	if err != nil {
		ts.Fatalf("AI error: %v", err)
	}
	if requests != 2 || !strings.Contains(string(loaded.RawContent()), "tone 220 binaural 40") {
		ts.Fatalf("unexpected semantic repair result: requests=%d content=%q", requests, loaded.RawContent())
	}
}

func TestAIStopsAfterTwoRepairs(ts *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"not valid spsq"}}]}`))
	}))
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	_, err := NewAppContext().AI(context.Background(), "make sequence", testAIOptions(server.URL))
	if err == nil || !strings.Contains(err.Error(), "after 2 repair attempts") {
		ts.Fatalf("unexpected error: %v", err)
	}
	if requests != 3 {
		ts.Fatalf("expected initial request plus two repairs, got %d", requests)
	}
}

func TestAIDoesNotRepairAPIErrors(ts *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(`{"error":{"message":"temporary failure"}}`))
	}))
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	_, err := NewAppContext().AI(context.Background(), "make sequence", testAIOptions(server.URL))
	if err == nil || !strings.Contains(err.Error(), "HTTP 500") {
		ts.Fatalf("unexpected error: %v", err)
	}
	if requests != 1 {
		ts.Fatalf("expected one API request, got %d", requests)
	}
}

func TestAIReportsCancellation(ts *testing.T) {
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewAppContext().AI(ctx, "make sequence", testAIOptions("http://127.0.0.1:1"))
	if err == nil || err.Error() != "AI generation canceled" {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestAIReportsTimeout(ts *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(50 * time.Millisecond)
	}))
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")
	timeout := 10 * time.Millisecond

	_, err := NewAppContext().AI(context.Background(), "make sequence", &AIOptions{
		BaseURL: server.URL,
		Timeout: timeout,
	})
	if err == nil || err.Error() != fmt.Sprintf("AI generation timed out after %s", timeout) {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestAIRejectsMissingTimeout(ts *testing.T) {
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	_, err := NewAppContext().AI(context.Background(), "make sequence", &AIOptions{BaseURL: "http://127.0.0.1:1"})
	if err == nil || err.Error() != "AI timeout must be greater than zero" {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func testAIOptions(baseURL string) *AIOptions {
	return &AIOptions{
		BaseURL:     baseURL,
		Temperature: 1,
		Timeout:     time.Minute,
	}
}
