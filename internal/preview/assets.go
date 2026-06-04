// Copyright (C) 2026 SynapSeq Contributors
//
// SPDX-License-Identifier: GPL-3.0-or-later

package preview

import (
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
)

//go:embed static/*
var staticFiles embed.FS

var previewTemplate = template.Must(template.ParseFS(staticFiles, "static/index.html"))

var previewCSS = template.CSS(mustReadPreviewAsset("static/preview.css"))

var previewJS = template.JS(mustReadPreviewAsset("static/preview.js"))

var lucideJS = template.JS(mustReadPreviewAsset("static/lucide.min.js"))

var chartJS = template.JS(mustReadPreviewAsset("static/chart.min.js"))

var previewLogoDataURL = template.URL(mustReadPreviewImageDataURL("static/logo.png"))

func mustReadPreviewAsset(path string) string {
	content, err := staticFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func mustReadPreviewImageDataURL(path string) string {
	content, err := staticFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(content))
}
