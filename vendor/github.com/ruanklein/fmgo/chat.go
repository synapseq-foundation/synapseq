// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"context"
	"fmt"
	"io"
	"sync"

	internalprocess "github.com/ruanklein/fmgo/internal/process"
)

// ChatOptions configures an interactive fm chat session.
type ChatOptions struct {
	Model           Model
	Instructions    string
	Resume          string
	Continue        bool
	SetDefaultModel Model
	Tools           []Tool
}

// ChatSession owns the standard streams and process for an interactive chat.
type ChatSession struct {
	command *internalprocess.PipedCommand
	args    []string
	once    sync.Once
}

// StartChat starts an interactive chat session.
func (c *Client) StartChat(ctx context.Context, options ChatOptions) (*ChatSession, error) {
	if options.Continue && options.Resume != "" {
		return nil, fmt.Errorf("fmgo: chat continue and resume cannot both be set")
	}
	args := []string{"chat"}
	if options.Model != "" {
		args = append(args, "--model", string(options.Model))
	}
	if options.Instructions != "" {
		args = append(args, "--instructions", options.Instructions)
	}
	if options.Resume != "" {
		args = append(args, "--resume", options.Resume)
	}
	if options.Continue {
		args = append(args, "--continue")
	}
	if options.SetDefaultModel != "" {
		args = append(args, "--set-default-model", string(options.SetDefaultModel))
	}
	for _, tool := range options.Tools {
		args = append(args, "--tool", string(tool))
	}
	path, err := c.executablePath()
	if err != nil {
		return nil, err
	}
	command, err := internalprocess.StartPiped(ctx, path, args...)
	if err != nil {
		return nil, commandError(args, "", err)
	}
	return &ChatSession{command: command, args: args}, nil
}

// Stdin returns the chat process input stream.
func (s *ChatSession) Stdin() io.WriteCloser { return s.command.Stdin() }

// Stdout returns the chat process output stream.
func (s *ChatSession) Stdout() io.ReadCloser { return s.command.Stdout() }

// Stderr returns the chat process error stream.
func (s *ChatSession) Stderr() io.ReadCloser { return s.command.Stderr() }

// Wait waits for the chat process after its output streams have been consumed.
func (s *ChatSession) Wait() error {
	if err := s.command.Wait(); err != nil {
		return commandError(s.args, "", err)
	}
	return nil
}

// Close terminates and reaps the chat process. It is safe to call more than once.
func (s *ChatSession) Close() error {
	var closeErr error
	s.once.Do(func() { closeErr = s.command.Close() })
	return closeErr
}
