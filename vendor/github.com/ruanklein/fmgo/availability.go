// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"context"
	"fmt"
	"strings"
)

// ModelAvailability describes whether a model is ready to use.
type ModelAvailability struct {
	Model     Model
	Available bool
	Reason    string
}

// Availability contains statuses reported by fm available.
type Availability struct {
	Models []ModelAvailability
}

// Available reports Foundation Models availability.
func (c *Client) Available(ctx context.Context) (Availability, error) {
	stdout, _, err := c.run(ctx, "available")
	if err != nil {
		return Availability{}, err
	}
	availability, err := parseAvailability(string(stdout))
	if err != nil {
		return Availability{}, err
	}
	return availability, nil
}

func parseAvailability(output string) (Availability, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	models := make([]ModelAvailability, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case line == "System model available":
			models = append(models, ModelAvailability{Model: ModelSystem, Available: true})
		case strings.HasPrefix(line, "System model unavailable:"):
			reason := strings.TrimSpace(strings.TrimPrefix(line, "System model unavailable:"))
			models = append(models, ModelAvailability{Model: ModelSystem, Reason: reason})
		default:
			return Availability{}, fmt.Errorf("fmgo: parse availability: %w", ErrUnexpectedOutput)
		}
	}
	if len(models) == 0 {
		return Availability{}, fmt.Errorf("fmgo: parse availability: %w", ErrUnexpectedOutput)
	}
	return Availability{Models: models}, nil
}
