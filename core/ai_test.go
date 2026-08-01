// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
