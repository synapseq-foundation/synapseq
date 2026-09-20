// SPDX-License-Identifier: Apache-2.0

package platform

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const MinimumMajorVersion = 27

var (
	// ErrUnsupportedPlatform reports a system other than macOS.
	ErrUnsupportedPlatform = errors.New("fmgo requires macOS")
	// ErrUnsupportedVersion reports a macOS version earlier than 27.
	ErrUnsupportedVersion = errors.New("fmgo requires macOS 27 or later")
)

// Check verifies that the current machine can run fmgo.
func Check() error {
	if runtime.GOOS != "darwin" {
		return ErrUnsupportedPlatform
	}

	output, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return fmt.Errorf("%w: unable to determine macOS product version", ErrUnsupportedVersion)
	}
	return Validate(runtime.GOOS, string(output))
}

// Validate verifies an operating system and macOS product version.
func Validate(goos, version string) error {
	if goos != "darwin" {
		return ErrUnsupportedPlatform
	}
	major, err := MajorVersion(version)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnsupportedVersion, err)
	}
	if major < MinimumMajorVersion {
		return fmt.Errorf("%w: found macOS %d", ErrUnsupportedVersion, major)
	}
	return nil
}

// MajorVersion extracts the major component of a macOS product version.
func MajorVersion(version string) (int, error) {
	majorText, _, _ := strings.Cut(strings.TrimSpace(version), ".")
	major, err := strconv.Atoi(majorText)
	if err != nil || major < 0 {
		return 0, fmt.Errorf("invalid macOS product version %q", version)
	}
	return major, nil
}
