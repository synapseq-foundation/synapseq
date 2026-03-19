/*
 * SynapSeq - Text-Driven Audio Sequencer for Brainwave Entrainment
 * https://synapseq.org
 *
 * Copyright (c) 2025-2026 SynapSeq Foundation
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2.
 * See the file COPYING.txt for details.
 */

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
