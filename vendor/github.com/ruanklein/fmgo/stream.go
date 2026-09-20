// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"context"
	"io"
	"sync"

	internalprocess "github.com/ruanklein/fmgo/internal/process"
)

// Stream incrementally reads a response from fm. Text returns the latest read block.
type Stream struct {
	command *internalprocess.PipedCommand
	stdout  io.ReadCloser
	stderr  <-chan []byte
	buffer  []byte
	text    string
	err     error
	once    sync.Once
	cleanup func()
	clean   sync.Once
}

// Stream starts a streaming response.
func (c *Client) Stream(ctx context.Context, request Request) (*Stream, error) {
	request, cleanup, err := request.withSchemaFile()
	if err != nil {
		return nil, err
	}
	args, err := request.args(true)
	if err != nil {
		cleanup()
		return nil, err
	}
	path, err := c.executablePath()
	if err != nil {
		cleanup()
		return nil, err
	}
	command, err := internalprocess.StartPiped(ctx, path, args...)
	if err != nil {
		cleanup()
		return nil, commandError(args, "", err)
	}
	stderr := make(chan []byte, 1)
	go func() {
		output, _ := io.ReadAll(command.Stderr())
		stderr <- output
	}()
	return &Stream{
		command: command,
		stdout:  command.Stdout(),
		stderr:  stderr,
		buffer:  make([]byte, 32*1024),
		cleanup: cleanup,
	}, nil
}

// Next reads the next available output block.
func (s *Stream) Next() bool {
	if s.err != nil {
		return false
	}
	for {
		n, err := s.stdout.Read(s.buffer)
		if n > 0 {
			s.text = string(s.buffer[:n])
			return true
		}
		if err == nil {
			continue
		}
		if err != io.EOF {
			s.err = err
			return false
		}
		break
	}
	stderr := <-s.stderr
	if waitErr := s.command.Wait(); waitErr != nil {
		s.err = commandError([]string{"respond", "--stream"}, string(stderr), waitErr)
	}
	s.cleanupSchemaFile()
	return false
}

// Text returns the text read by the most recent successful Next call.
func (s *Stream) Text() string { return s.text }

// Err returns the terminal stream error, if any.
func (s *Stream) Err() error { return s.err }

// Close terminates the response process. It is safe to call more than once.
func (s *Stream) Close() error {
	var closeErr error
	s.once.Do(func() {
		_ = s.stdout.Close()
		closeErr = s.command.Close()
		s.cleanupSchemaFile()
	})
	return closeErr
}

func (s *Stream) cleanupSchemaFile() {
	s.clean.Do(s.cleanup)
}
