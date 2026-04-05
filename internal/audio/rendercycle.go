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
	"math"

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
		totalFrames: totalRenderFrames(renderer.periods, renderer.SampleRate),
		chunkFrames: int64(t.BufferSize),
	}

	if renderer.StatusOutput != nil {
		runtime.status = audiostatus.NewReporter(renderer.StatusOutput, renderer.Colors)
	}

	return runtime
}

func totalRenderFrames(periods []t.Period, sampleRate int) int64 {
	endMs := periods[len(periods)-1].Time
	return int64(math.Round(float64(endMs) * float64(sampleRate) / 1000.0))
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
	for rr.periodIdx+1 < len(rr.renderer.periods) && currentTimeMs >= rr.renderer.periods[rr.periodIdx+1].Time {
		rr.periodIdx++
	}
}

func (rr *renderRuntime) syncAndPrepare(currentTimeMs int) {
	rr.renderer.syncEngine.Sync(rr.renderer.channels[:], rr.renderer.periods, currentTimeMs, rr.periodIdx)
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
