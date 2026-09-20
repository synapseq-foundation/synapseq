// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// TokenRequest describes content to count with fm count-tokens.
type TokenRequest struct {
	Prompt       string
	Instructions string
	Text         []string
	Images       []string
	Transcript   string
}

// CountTokens counts Foundation Models tokens.
func (c *Client) CountTokens(ctx context.Context, request TokenRequest) (int, error) {
	args := []string{"count-tokens", "--quiet"}
	if request.Instructions != "" {
		args = append(args, "--instructions", request.Instructions)
	}
	for _, image := range request.Images {
		args = append(args, "--image", image)
	}
	for _, text := range request.Text {
		args = append(args, "--text", text)
	}
	if request.Transcript != "" {
		args = append(args, "--transcript", request.Transcript)
	}
	if request.Prompt != "" {
		args = append(args, request.Prompt)
	}
	stdout, _, err := c.run(ctx, args...)
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
	if err != nil {
		return 0, fmt.Errorf("fmgo: parse token count: %w", ErrUnexpectedOutput)
	}
	return count, nil
}
