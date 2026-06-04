// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package parser

import (
	"fmt"
	"strings"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

// HasComment checks if the first element is a comment
func (ctx *TextParser) HasComment() bool {
	tok, ok := ctx.Line.Peek()
	return ok && string(tok[0]) == t.KeywordComment
}

// ParseComment extracts and prints the comment from the elements
func (ctx *TextParser) ParseComment() string {
	tok, ok := ctx.Line.Peek()
	if !ok || string(tok[0]) != t.KeywordComment {
		return ""
	}
	if len(tok) >= 2 && string(tok[1]) == t.KeywordComment {
		comment := fmt.Sprintf("%s %s", tok[2:], strings.Join(ctx.Line.Tokens[1:], " "))
		// Trim leading/trailing whitespace if there's more than just the ##
		if len(comment) > 1 {
			comment = strings.TrimSpace(comment)
		}
		return comment
	}
	return ""
}
