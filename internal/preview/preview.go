// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

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

	sectionData := buildPreviewSectionData(periods, totalDurationMs)

	return buildPreviewTemplateData(periods, totalDurationMs, sectionData), nil
}
