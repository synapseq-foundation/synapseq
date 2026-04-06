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

package preview

import (
	"bytes"

	t "github.com/synapseq-foundation/synapseq/v4/internal/types"
)

func GetPreviewContent(periods []t.Period) ([]byte, error) {
	data, err := buildPreviewData(periods)
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	if err := previewTemplate.Execute(&output, data); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func buildPreviewData(periods []t.Period) (*previewTemplateData, error) {
	totalDurationMs, err := validatePreviewPeriods(periods)
	if err != nil {
		return nil, err
	}

	renderData := buildPreviewRenderData(periods, totalDurationMs)

	return buildPreviewTemplateData(periods, totalDurationMs, renderData), nil
}
