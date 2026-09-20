// SPDX-License-Identifier: Apache-2.0

package process

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"sync"
)

// Result contains output collected from a completed command.
type Result struct {
	Stdout []byte
	Stderr []byte
}

// Run executes a command without a shell.
func Run(ctx context.Context, executable string, args ...string) (Result, error) {
	cmd := exec.CommandContext(ctx, executable, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return Result{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}, err
}

// PipedCommand is a started command with its standard streams exposed.
type PipedCommand struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	wait   sync.Once
	err    error
}

// StartPiped starts a command and exposes all of its standard streams.
func StartPiped(ctx context.Context, executable string, args ...string) (*PipedCommand, error) {
	cmd := exec.CommandContext(ctx, executable, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &PipedCommand{cmd: cmd, stdin: stdin, stdout: stdout, stderr: stderr}, nil
}

func (c *PipedCommand) Stdin() io.WriteCloser { return c.stdin }
func (c *PipedCommand) Stdout() io.ReadCloser { return c.stdout }
func (c *PipedCommand) Stderr() io.ReadCloser { return c.stderr }

// Wait waits for the command exactly once.
func (c *PipedCommand) Wait() error {
	c.wait.Do(func() { c.err = c.cmd.Wait() })
	return c.err
}

// Close terminates the command and reaps it.
func (c *PipedCommand) Close() error {
	_ = c.stdin.Close()
	if c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
	}
	_ = c.Wait()
	return nil
}

// ManagedCommand is a started command whose output is retained for diagnostics.
type ManagedCommand struct {
	cmd    *exec.Cmd
	stderr bytes.Buffer
	err    error
	done   chan struct{}
}

// Start starts a command and retains standard error.
func Start(ctx context.Context, executable string, args ...string) (*ManagedCommand, error) {
	cmd := exec.CommandContext(ctx, executable, args...)
	managed := &ManagedCommand{cmd: cmd, done: make(chan struct{})}
	cmd.Stderr = &managed.stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	go func() {
		managed.err = cmd.Wait()
		close(managed.done)
	}()
	return managed, nil
}

func (c *ManagedCommand) Stderr() string { return c.stderr.String() }

// Wait waits for the command exactly once.
func (c *ManagedCommand) Wait() error {
	<-c.done
	return c.err
}

// Done closes when the command exits.
func (c *ManagedCommand) Done() <-chan struct{} { return c.done }

// Close terminates the command and reaps it.
func (c *ManagedCommand) Close() error {
	if c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
	}
	_ = c.Wait()
	return nil
}
