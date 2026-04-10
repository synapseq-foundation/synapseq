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

package audio

import (
	"fmt"

	audiostatus "github.com/synapseq-foundation/synapseq/v4/internal/audio/status"
	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

type renderRuntime struct {
	renderer      *AudioRenderer
	consume       func(samples []int) error
	status        *audiostatus.Reporter
	samples       []int
	totalFrames   int64
	chunkFrames   int64
	framesWritten int64
	periodIdx     int
}

func newRenderRuntime(renderer *AudioRenderer, consume func(samples []int) error) *renderRuntime {
	runtime := &renderRuntime{
		renderer:    renderer,
		consume:     consume,
		samples:     make([]int, t.BufferSize*audioChannels),
		totalFrames: renderer.plan.totalFrames,
		chunkFrames: int64(t.BufferSize),
	}

	if renderer.StatusOutput != nil {
		runtime.status = audiostatus.NewReporter(renderer.StatusOutput, renderer.Colors)
	}

	return runtime
}

func (rr *renderRuntime) run() error {
	if rr.status != nil {
		defer rr.status.FinalStatus()
	}

	for rr.framesWritten < rr.totalFrames {
		currentTimeMs := rr.currentTimeMs()
		rr.advancePeriod(currentTimeMs)
		rr.syncAndPrepare(currentTimeMs)
		rr.reportPeriodChange()

		data, framesToWrite := rr.mixCurrentChunk()
		if err := rr.consumeChunk(data); err != nil {
			return err
		}

		rr.framesWritten += framesToWrite
		rr.reportProgress(currentTimeMs)
	}

	return nil
}

func (rr *renderRuntime) currentTimeMs() int {
	return int((float64(rr.framesWritten) * 1000.0) / float64(rr.renderer.SampleRate))
}

func (rr *renderRuntime) advancePeriod(currentTimeMs int) {
	rr.periodIdx = rr.renderer.plan.periodIndexAt(currentTimeMs, rr.periodIdx)
}

func (rr *renderRuntime) syncAndPrepare(currentTimeMs int) {
	cue := rr.renderer.plan.cue(rr.periodIdx, currentTimeMs)
	rr.renderer.applyCueSignalState(cue)
	rr.renderer.syncEngine.Sync(rr.renderer.channels[:], cue)
	if rr.renderer.ambianceState != nil {
		rr.renderer.ambianceState.CollectActiveIndices(rr.renderer.channels[:])
		rr.renderer.ambianceState.PrepareBuffers(t.BufferSize)
	}
}

func (rr *renderRuntime) reportPeriodChange() {
	if rr.status != nil {
		rr.status.CheckPeriodChange(rr.statusView(), rr.periodIdx)
	}
}

func (rr *renderRuntime) mixCurrentChunk() ([]int, int64) {
	data := rr.renderer.mix(rr.samples)
	framesToWrite := rr.chunkFrames

	if remain := rr.totalFrames - rr.framesWritten; remain < rr.chunkFrames {
		framesToWrite = remain
		data = data[:remain*audioChannels]
	}

	return data, framesToWrite
}

func (rr *renderRuntime) consumeChunk(data []int) error {
	if rr.consume == nil {
		return nil
	}

	if err := rr.consume(data); err != nil {
		return fmt.Errorf("failed to consume audio buffer: %w", err)
	}

	return nil
}

func (rr *renderRuntime) reportProgress(currentTimeMs int) {
	if rr.status != nil && rr.status.ShouldUpdateStatus() {
		rr.status.DisplayStatus(rr.statusView(), currentTimeMs)
	}
}

func (rr *renderRuntime) statusView() audiostatus.View {
	return audiostatus.View{
		Periods:  rr.renderer.periods,
		Channels: rr.renderer.channels[:],
	}
}
