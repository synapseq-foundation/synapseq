// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

//go:build darwin && arm64

package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/ruanklein/fmgo"
)

type foundationClient struct {
	client *fmgo.Client
}

func newFoundationClient() (foundationClient, error) {
	client, err := fmgo.New()
	if err != nil {
		return foundationClient{}, err
	}
	return foundationClient{client: client}, nil
}

func (c foundationClient) generate(ctx context.Context, config Config, prompt string) (string, error) {
	response, err := c.client.Respond(ctx, fmgo.Request{
		Prompt:       prompt,
		Instructions: appleFoundationSystemPromptForRequest(prompt),
		Model:        fmgo.Model(config.Model),
	})
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(response.Text) == "" {
		return "", fmt.Errorf("AI did not understand the prompt")
	}

	return strings.TrimSpace(response.Text) + "\n", nil
}
