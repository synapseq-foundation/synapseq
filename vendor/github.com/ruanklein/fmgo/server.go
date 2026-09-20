// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	internalprocess "github.com/ruanklein/fmgo/internal/process"
)

// ServerOptions configures the native fm serve process.
type ServerOptions struct {
	Host           string
	Port           int
	Socket         string
	StartupTimeout time.Duration
}

// Server manages a native fm serve process.
type Server struct {
	command *internalprocess.ManagedCommand
	address string
	once    sync.Once
}

// Serve starts fm serve and waits until its configured address accepts connections.
func (c *Client) Serve(ctx context.Context, options ServerOptions) (*Server, error) {
	args, network, address, err := serverArgs(options)
	if err != nil {
		return nil, err
	}
	path, err := c.executablePath()
	if err != nil {
		return nil, err
	}
	command, err := internalprocess.Start(ctx, path, args...)
	if err != nil {
		return nil, commandError(args, "", err)
	}
	server := &Server{command: command, address: address}
	if err := waitForAddress(ctx, network, address, options.StartupTimeout, command); err != nil {
		_ = command.Close()
		return nil, err
	}
	return server, nil
}

// Addr returns the server's configured network address.
func (s *Server) Addr() string { return s.address }

// Wait waits for the native server to exit.
func (s *Server) Wait() error {
	err := s.command.Wait()
	if err != nil {
		return commandError([]string{"serve"}, s.command.Stderr(), err)
	}
	return nil
}

// Close terminates and reaps the native server. It is safe to call more than once.
func (s *Server) Close() error {
	var closeErr error
	s.once.Do(func() {
		closeErr = s.command.Close()
	})
	return closeErr
}

func serverArgs(options ServerOptions) ([]string, string, string, error) {
	if options.Socket != "" {
		if options.Host != "" || options.Port != 0 {
			return nil, "", "", fmt.Errorf("fmgo: socket cannot be combined with host or port")
		}
		return []string{"serve", "--socket", options.Socket}, "unix", options.Socket, nil
	}
	if options.Port < 1 || options.Port > 65535 {
		return nil, "", "", fmt.Errorf("fmgo: port must be between 1 and 65535")
	}
	host := options.Host
	if host == "" {
		host = "127.0.0.1"
	}
	address := net.JoinHostPort(host, fmt.Sprintf("%d", options.Port))
	return []string{"serve", "--host", host, "--port", fmt.Sprintf("%d", options.Port)}, "tcp", address, nil
}

func waitForAddress(
	ctx context.Context,
	network string,
	address string,
	timeout time.Duration,
	command *internalprocess.ManagedCommand,
) error {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	startupContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	dialer := net.Dialer{}
	for {
		connection, err := dialer.DialContext(startupContext, network, address)
		if err == nil {
			return connection.Close()
		}
		select {
		case <-command.Done():
			return commandError([]string{"serve"}, command.Stderr(), command.Wait())
		case <-startupContext.Done():
			return fmt.Errorf("fmgo: server startup timed out after %s", timeout)
		case <-time.After(25 * time.Millisecond):
		}
	}
}
