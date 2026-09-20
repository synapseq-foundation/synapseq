// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"context"

	internalprocess "github.com/ruanklein/fmgo/internal/process"
)

type commandResult struct {
	stdout []byte
	stderr string
}

func runCommand(ctx context.Context, executable string, args ...string) (commandResult, error) {
	result, err := internalprocess.Run(ctx, executable, args...)
	return commandResult{stdout: result.Stdout, stderr: string(result.Stderr)}, err
}
