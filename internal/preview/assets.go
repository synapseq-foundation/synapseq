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
	"embed"
	"html/template"
)

//go:embed static/*
var staticFiles embed.FS

var previewTemplate = template.Must(template.ParseFS(staticFiles, "static/index.html"))

var previewCSS = template.CSS(mustReadPreviewAsset("static/preview.css"))

var previewJS = template.JS(mustReadPreviewAsset("static/preview.js"))

func mustReadPreviewAsset(path string) string {
	content, err := staticFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(content)
}
