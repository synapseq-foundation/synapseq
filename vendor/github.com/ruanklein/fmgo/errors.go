// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ruanklein/fmgo/internal/platform"
)

var (
	// ErrUnsupportedPlatform reports a system other than macOS.
	ErrUnsupportedPlatform = platform.ErrUnsupportedPlatform
	// ErrUnsupportedVersion reports a macOS version earlier than 27.
	ErrUnsupportedVersion = platform.ErrUnsupportedVersion
	// ErrFMNotFound reports that the native fm executable could not be located.
	ErrFMNotFound = errors.New("fm executable not found")
	// ErrModelUnavailable reports that Foundation Models are not ready on this machine.
	ErrModelUnavailable = errors.New("foundation model unavailable")
	// ErrLicenseRequired reports that the Foundation Models CLI terms have not been accepted.
	ErrLicenseRequired = errors.New("foundation models license not accepted")
	// ErrUnexpectedOutput reports output that fmgo cannot safely interpret.
	ErrUnexpectedOutput = errors.New("unexpected fm output")
)

// CommandError describes a failed native fm invocation.
type CommandError struct {
	Command  string
	ExitCode int
	Stderr   string
	cause    error
}

func (e *CommandError) Error() string {
	if e.ExitCode >= 0 {
		return fmt.Sprintf("fmgo: %s exited with status %d", e.Command, e.ExitCode)
	}
	return fmt.Sprintf("fmgo: %s failed", e.Command)
}

// Unwrap returns a recognized operational cause when one is available.
func (e *CommandError) Unwrap() error { return e.cause }

func commandError(args []string, stderr string, err error) error {
	if isNotFound(err) {
		return ErrFMNotFound
	}

	exitCode := -1
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		exitCode = exitError.ExitCode()
	}

	cause := error(nil)
	lowerStderr := strings.ToLower(stderr)
	switch {
	case exitCode == 69:
		cause = ErrLicenseRequired
	case strings.Contains(lowerStderr, "model is not available") || strings.Contains(lowerStderr, "modelnotready"):
		cause = ErrModelUnavailable
	}

	return &CommandError{
		Command:  strings.Join(append([]string{"fm"}, args...), " "),
		ExitCode: exitCode,
		Stderr:   stderr,
		cause:    cause,
	}
}
