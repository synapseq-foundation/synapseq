// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"context"
	"fmt"
	"strings"
)

// LicenseStatus reports whether the local machine accepted the fm terms.
type LicenseStatus struct {
	Accepted bool
}

// License reports whether the Foundation Models CLI terms were accepted.
func (c *Client) License(ctx context.Context) (LicenseStatus, error) {
	stdout, _, err := c.run(ctx, "license", "--status")
	if err != nil {
		return LicenseStatus{}, err
	}
	if strings.HasPrefix(strings.TrimSpace(string(stdout)), "Agreed to license ") {
		return LicenseStatus{Accepted: true}, nil
	}
	return LicenseStatus{}, fmt.Errorf("fmgo: parse license status: %w", ErrUnexpectedOutput)
}

// ShowLicense returns the Foundation Models CLI terms without prompting.
func (c *Client) ShowLicense(ctx context.Context) (string, error) {
	stdout, _, err := c.run(ctx, "license", "--show")
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}
