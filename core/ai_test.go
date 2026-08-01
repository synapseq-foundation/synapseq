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

	loaded, err := NewAppContext().AI("make relaxation", &AIOptions{BaseURL: server.URL})
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

	_, err := NewAppContext().AI("make relaxation", &AIOptions{BaseURL: server.URL})
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

func TestAITemperatureUsesOptionsBeforeEnvironment(ts *testing.T) {
	ts.Setenv("SYNAPSEQ_AI_TEMPERATURE", "0.2")
	temperature := 0.7

	got, err := aiTemperature(&AIOptions{Temperature: &temperature})
	if err != nil {
		ts.Fatalf("aiTemperature error: %v", err)
	}
	if got == nil || *got != 0.7 {
		ts.Fatalf("expected option temperature 0.7, got %#v", got)
	}
}

func TestAITemperatureUsesEnvironment(ts *testing.T) {
	ts.Setenv("SYNAPSEQ_AI_TEMPERATURE", "0.2")

	got, err := aiTemperature(nil)
	if err != nil {
		ts.Fatalf("aiTemperature error: %v", err)
	}
	if got == nil || *got != 0.2 {
		ts.Fatalf("expected environment temperature 0.2, got %#v", got)
	}
}

func TestAITemperatureRejectsInvalidValue(ts *testing.T) {
	ts.Setenv("SYNAPSEQ_AI_TEMPERATURE", "3")

	if _, err := aiTemperature(nil); err == nil {
		ts.Fatal("expected invalid temperature error")
	}
}

func TestAIContextRepairsInvalidSPSQ(ts *testing.T) {
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

	loaded, err := NewAppContext().AIContext(context.Background(), "make relaxation", &AIOptions{BaseURL: server.URL})
	if err != nil {
		ts.Fatalf("AIContext error: %v", err)
	}
	if requests != 2 || !strings.Contains(string(loaded.RawContent()), "relax") {
		ts.Fatalf("unexpected repair result: requests=%d content=%q", requests, loaded.RawContent())
	}
}

func TestAIContextStopsAfterTwoRepairs(ts *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		_, _ = writer.Write([]byte(`{"choices":[{"message":{"content":"not valid spsq"}}]}`))
	}))
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	_, err := NewAppContext().AIContext(context.Background(), "make relaxation", &AIOptions{BaseURL: server.URL})
	if err == nil || !strings.Contains(err.Error(), "after 2 repair attempts") {
		ts.Fatalf("unexpected error: %v", err)
	}
	if requests != 3 {
		ts.Fatalf("expected initial request plus two repairs, got %d", requests)
	}
}

func TestAIContextDoesNotRepairAPIErrors(ts *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(`{"error":{"message":"temporary failure"}}`))
	}))
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")

	_, err := NewAppContext().AIContext(context.Background(), "make relaxation", &AIOptions{BaseURL: server.URL})
	if err == nil || !strings.Contains(err.Error(), "HTTP 500") {
		ts.Fatalf("unexpected error: %v", err)
	}
	if requests != 1 {
		ts.Fatalf("expected one API request, got %d", requests)
	}
}

func TestAIContextReportsCancellation(ts *testing.T) {
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewAppContext().AIContext(ctx, "make relaxation", &AIOptions{BaseURL: "http://127.0.0.1:1"})
	if err == nil || err.Error() != "AI generation canceled" {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestAIContextReportsTimeout(ts *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(50 * time.Millisecond)
	}))
	defer server.Close()
	ts.Setenv("SYNAPSEQ_AI_API_KEY", "test-key")
	timeout := 10 * time.Millisecond

	_, err := NewAppContext().AIContext(context.Background(), "make relaxation", &AIOptions{
		BaseURL: server.URL,
		Timeout: &timeout,
	})
	if err == nil || err.Error() != fmt.Sprintf("AI generation timed out after %s", timeout) {
		ts.Fatalf("unexpected error: %v", err)
	}
}

func TestAITimeoutUsesOptionsBeforeEnvironment(ts *testing.T) {
	ts.Setenv("SYNAPSEQ_AI_TIMEOUT", "2m")
	timeout := 30 * time.Second

	got, err := aiTimeout(&AIOptions{Timeout: &timeout})
	if err != nil {
		ts.Fatalf("aiTimeout error: %v", err)
	}
	if got != timeout {
		ts.Fatalf("expected option timeout %s, got %s", timeout, got)
	}
}

func TestAITimeoutUsesDefault(ts *testing.T) {
	ts.Setenv("SYNAPSEQ_AI_TIMEOUT", "")

	got, err := aiTimeout(nil)
	if err != nil {
		ts.Fatalf("aiTimeout error: %v", err)
	}
	if got != defaultAITimeout {
		ts.Fatalf("expected default timeout %s, got %s", defaultAITimeout, got)
	}
}
