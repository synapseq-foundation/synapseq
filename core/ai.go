// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package core

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	internalai "github.com/synapseq-foundation/synapseq/v4/internal/ai"
)

const (
	defaultAIModel   = "gpt-4.1-mini"
	aiRepairAttempts = 2
)

// AI generates and validates an SPSQ sequence from prompt using an
// OpenAI-compatible chat completion API. The API key is read from
// SYNAPSEQ_AI_API_KEY. The request stops when ctx is canceled or its configured
// timeout expires.
func (ac *AppContext) AI(ctx context.Context, prompt string, options *AIOptions) (*LoadedContext, error) {
	if ac == nil {
		return nil, fmt.Errorf("app context is nil")
	}
	if ctx == nil {
		return nil, fmt.Errorf("AI context is nil")
	}
	if options == nil {
		return nil, fmt.Errorf("AI options are nil")
	}

	if err := validateAITemperature(options.Temperature); err != nil {
		return nil, err
	}
	timeout, err := validateAITimeout(options.Timeout)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := internalai.New(internalai.Config{
		APIKey:      os.Getenv("SYNAPSEQ_AI_API_KEY"),
		BaseURL:     aiBaseURL(options),
		Model:       aiModel(options),
		Temperature: options.Temperature,
	})
	if err != nil {
		return nil, err
	}

	content, err := client.Generate(ctx, prompt)
	if err != nil {
		return nil, aiContextError(ctx, timeout, err)
	}

	for attempt := 0; attempt <= aiRepairAttempts; attempt++ {
		loaded, validationErr := ac.LoadContent(content)
		if validationErr == nil {
			validationErr = validateAISequenceSemantics(loaded, prompt)
			if validationErr == nil {
				return loaded, nil
			}
		}
		if attempt == aiRepairAttempts {
			return nil, fmt.Errorf("AI did not understand the prompt after %d repair attempts: generated content is not valid SPSQ: %w", aiRepairAttempts, validationErr)
		}

		content, err = client.Repair(ctx, prompt, content, validationErr)
		if err != nil {
			return nil, aiContextError(ctx, timeout, err)
		}
	}

	return nil, fmt.Errorf("AI generation ended unexpectedly")
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

func validateAITemperature(temperature float64) error {
	if math.IsNaN(temperature) || math.IsInf(temperature, 0) || temperature < 0 || temperature > 2 {
		return fmt.Errorf("AI temperature must be between 0 and 2")
	}

	return nil
}

func validateAITimeout(timeout time.Duration) (time.Duration, error) {
	if timeout <= 0 {
		return 0, fmt.Errorf("AI timeout must be greater than zero")
	}

	return timeout, nil
}

func aiContextError(ctx context.Context, timeout time.Duration, err error) error {
	if errors.Is(ctx.Err(), context.Canceled) {
		return fmt.Errorf("AI generation canceled")
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("AI generation timed out after %s", timeout)
	}

	return err
}
