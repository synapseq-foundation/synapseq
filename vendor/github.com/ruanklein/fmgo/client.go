// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"context"
	"errors"
	"os/exec"

	"github.com/ruanklein/fmgo/internal/platform"
)

// Client invokes the native fm executable.
type Client struct {
	executable string
}

// Option configures a Client.
type Option func(*Client)

// WithExecutable overrides the fm executable path.
func WithExecutable(path string) Option {
	return func(client *Client) {
		client.executable = path
	}
}

// New creates a Client that uses fm from PATH by default.
func New(options ...Option) (*Client, error) {
	if err := platform.Check(); err != nil {
		return nil, err
	}
	client := &Client{executable: "fm"}
	for _, option := range options {
		option(client)
	}
	return client, nil
}

func (c *Client) executablePath() (string, error) {
	path, err := exec.LookPath(c.executable)
	if err != nil {
		return "", ErrFMNotFound
	}
	return path, nil
}

func (c *Client) run(ctx context.Context, args ...string) ([]byte, string, error) {
	path, err := c.executablePath()
	if err != nil {
		return nil, "", err
	}
	result, err := runCommand(ctx, path, args...)
	if err != nil {
		return result.stdout, result.stderr, commandError(args, result.stderr, err)
	}
	return result.stdout, result.stderr, nil
}

func isNotFound(err error) bool {
	return errors.Is(err, exec.ErrNotFound)
}
