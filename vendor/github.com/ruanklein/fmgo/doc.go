// SPDX-License-Identifier: Apache-2.0

// Package fmgo provides pure Go access to Apple's Foundation Models CLI.
//
// fmgo invokes the native fm executable. It is not a Swift bridge, does not
// require cgo, does not access private Apple frameworks, and does not
// reimplement Foundation Models or provide its own CLI.
//
// fmgo requires macOS 27 or later, the native fm executable, an available
// Foundation Model, and accepted Foundation Models CLI terms. New validates the
// platform and returns ErrUnsupportedPlatform or ErrUnsupportedVersion when the
// machine is incompatible. Operations return ErrFMNotFound when the executable
// cannot be located.
//
// Create a Client with New, then use Respond or Stream for generation.
// SchemaFor and RespondAs support structured output. Client also exposes token
// counting, availability and terms inspection, interactive chat sessions, and
// lifecycle management for the native fm serve process. All blocking operations
// accept context.Context.
package fmgo
