// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
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

	temperature, err := aiTemperature(options)
	if err != nil {
		return nil, err
	}

	client, err := internalai.New(internalai.Config{
		APIKey:      os.Getenv("SYNAPSEQ_AI_API_KEY"),
		BaseURL:     aiBaseURL(options),
		Model:       aiModel(options),
		Temperature: temperature,
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

func aiTemperature(options *AIOptions) (*float64, error) {
	if options != nil && options.Temperature != nil {
		return options.Temperature, validateAITemperature(*options.Temperature)
	}

	value := strings.TrimSpace(os.Getenv("SYNAPSEQ_AI_TEMPERATURE"))
	if value == "" {
		return nil, nil
	}

	temperature, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid SYNAPSEQ_AI_TEMPERATURE %q: %w", value, err)
	}

	return &temperature, validateAITemperature(temperature)
}

func validateAITemperature(temperature float64) error {
	if math.IsNaN(temperature) || math.IsInf(temperature, 0) || temperature < 0 || temperature > 2 {
		return fmt.Errorf("AI temperature must be between 0 and 2")
	}

	return nil
}
