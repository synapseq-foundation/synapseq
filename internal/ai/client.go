// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

// Package ai communicates with OpenAI-compatible chat completion APIs.
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const defaultBaseURL = "https://api.openai.com/v1"

type Config struct {
	APIKey      string
	BaseURL     string
	Model       string
	Temperature *float64
}

type Client struct {
	config     Config
	httpClient *http.Client
}

type chatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature *float64  `json:"temperature,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message message `json:"message"`
}

type apiErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func New(config Config) (*Client, error) {
	if strings.TrimSpace(config.APIKey) == "" {
		return nil, fmt.Errorf("SYNAPSEQ_AI_API_KEY is not set")
	}
	if strings.TrimSpace(config.Model) == "" {
		return nil, fmt.Errorf("AI model cannot be empty")
	}

	baseURL, err := completionURL(config.BaseURL)
	if err != nil {
		return nil, err
	}
	config.BaseURL = baseURL

	return &Client{
		config:     config,
		httpClient: http.DefaultClient,
	}, nil
}

func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	if strings.TrimSpace(prompt) == "" {
		return "", fmt.Errorf("AI prompt cannot be empty")
	}

	body, err := json.Marshal(chatCompletionRequest{
		Model: c.config.Model,
		Messages: []message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		Temperature: c.config.Temperature,
	})
	if err != nil {
		return "", fmt.Errorf("encode AI request: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.config.BaseURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("create AI request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("request AI completion: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read AI response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", apiError(response.StatusCode, responseBody)
	}

	var completion chatCompletionResponse
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return "", fmt.Errorf("decode AI response: %w", err)
	}
	if len(completion.Choices) == 0 || strings.TrimSpace(completion.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("AI did not understand the prompt")
	}

	return strings.TrimSpace(completion.Choices[0].Message.Content) + "\n", nil
}

func completionURL(baseURL string) (string, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid AI base URL %q", baseURL)
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	if !strings.HasSuffix(parsed.Path, "/v1") {
		parsed.Path += "/v1"
	}
	parsed.Path += "/chat/completions"

	return parsed.String(), nil
}

func apiError(statusCode int, body []byte) error {
	var response apiErrorResponse
	if err := json.Unmarshal(body, &response); err == nil && response.Error.Message != "" {
		return fmt.Errorf("AI API returned HTTP %d: %s", statusCode, response.Error.Message)
	}

	return fmt.Errorf("AI API returned HTTP %d", statusCode)
}
