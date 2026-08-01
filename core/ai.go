// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"context"
	"fmt"
	"os"
	"strings"

	internalai "github.com/synapseq-foundation/synapseq/v4/internal/ai"
)

const defaultAIModel = "gpt-4.1-mini"

// AI generates and validates an SPSQ sequence from prompt using an
// OpenAI-compatible chat completion API. The API key is read from
// SYNAPSEQ_AI_API_KEY.
func (ac *AppContext) AI(prompt string, options *AIOptions) (*LoadedContext, error) {
	if ac == nil {
		return nil, fmt.Errorf("app context is nil")
	}

	client, err := internalai.New(internalai.Config{
		APIKey:  os.Getenv("SYNAPSEQ_AI_API_KEY"),
		BaseURL: aiBaseURL(options),
		Model:   aiModel(options),
	})
	if err != nil {
		return nil, err
	}

	content, err := client.Generate(context.Background(), prompt)
	if err != nil {
		return nil, err
	}

	loaded, err := ac.LoadContent(content)
	if err != nil {
		return nil, fmt.Errorf("AI did not understand the prompt: generated content is not valid SPSQ: %w", err)
	}

	return loaded, nil
}

func aiModel(options *AIOptions) string {
	if options != nil && strings.TrimSpace(options.Model) != "" {
		return options.Model
	}
	if model := strings.TrimSpace(os.Getenv("SYNAPSEQ_AI_MODEL")); model != "" {
		return model
	}

	return defaultAIModel
}

func aiBaseURL(options *AIOptions) string {
	if options != nil && strings.TrimSpace(options.BaseURL) != "" {
		return options.BaseURL
	}

	return os.Getenv("SYNAPSEQ_AI_BASE_URL")
}
