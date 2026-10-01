// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/synapseq-foundation/synapseq/v4/internal/cli"
)

const aiProgressMessage = "Generating SPSQ sequence with AI. Please wait..."

type aiProgress struct {
	stop chan struct{}
	done chan struct{}
}

func startAIProgress(writer io.Writer, quiet bool) *aiProgress {
	if quiet || writer == nil {
		return nil
	}
	if !canAnimateAIProgress(writer) {
		_, _ = fmt.Fprintln(writer, cli.Muted(aiProgressMessage))
		return nil
	}

	progress := &aiProgress{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go progress.animate(writer)

	return progress
}

func (p *aiProgress) Stop() {
	if p == nil {
		return
	}

	close(p.stop)
	<-p.done
}

func (p *aiProgress) animate(writer io.Writer) {
	defer close(p.done)

	frames := []string{"|", "/", "-", `\`}
	writeFrame := func(index int) {
		_, _ = fmt.Fprintf(writer, "\r%s %s", cli.Accent(frames[index]), cli.Muted(aiProgressMessage))
	}
	writeFrame(0)

	ticker := time.NewTicker(120 * time.Millisecond)
	defer ticker.Stop()

	for index := 1; ; index = (index + 1) % len(frames) {
		select {
		case <-p.stop:
			_, _ = fmt.Fprintln(writer)
			return
		case <-ticker.C:
			writeFrame(index)
		}
	}
}

func canAnimateAIProgress(writer io.Writer) bool {
	file, ok := writer.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}
