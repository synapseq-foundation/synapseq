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

package sequence

import (
	"bufio"
	"bytes"
)

// SequenceFile represents a sequence file
type SequenceFile struct {
	currentLine       string         // Current line in the file
	currentLineNumber int            // Current line number
	scanner           *bufio.Scanner // Scanner for reading the file
}

// NewSequenceFile creates a new sequence file
func NewSequenceFile(data []byte) *SequenceFile {
	return &SequenceFile{
		scanner: bufio.NewScanner(bytes.NewReader(data)),
	}
}

// NextLine advances to the next line in the sequence file
func (sf *SequenceFile) NextLine() bool {
	if sf.scanner == nil {
		return false
	}

	if sf.scanner.Scan() {
		sf.currentLine = sf.scanner.Text()
		sf.currentLineNumber++
		return true
	}
	return false
}

// CurrentLine returns the current line in the sequence file
func (sf *SequenceFile) CurrentLine() string {
	return sf.currentLine
}

// CurrentLineNumber returns the current line number in the sequence file
func (sf *SequenceFile) CurrentLineNumber() int {
	return sf.currentLineNumber
}
