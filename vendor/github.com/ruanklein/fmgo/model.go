// SPDX-License-Identifier: Apache-2.0

package fmgo

// Model identifies a Foundation Model exposed by fm.
type Model string

const (
	// ModelSystem is the on-device Apple Foundation Model.
	ModelSystem Model = "system"
)

// UseCase configures the system model's intended task.
type UseCase string

const (
	UseCaseGeneral        UseCase = "general"
	UseCaseContentTagging UseCase = "content-tagging"
)

// GuardrailLevel configures Apple's supported guardrail level.
type GuardrailLevel string

const (
	GuardrailDefault                          GuardrailLevel = "default"
	GuardrailPermissiveContentTransformations GuardrailLevel = "permissive-content-transformations"
)

// Tool identifies a native fm tool.
type Tool string

const (
	ToolBarcode Tool = "barcode"
	ToolOCR     Tool = "ocr"
)
