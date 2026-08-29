# Go API

SynapSeq exposes a public Go API for generating, validating, converting, and
rendering sequences programmatically. The public packages are documented on
[`pkg.go.dev`](https://pkg.go.dev/github.com/synapseq-foundation/synapseq/v4).

## AI Generation From Go

`AppContext.AI` generates and validates SPSQ through an OpenAI-compatible API.
Set `SYNAPSEQ_AI_API_KEY` before calling it. `AIOptions` can select a local
model, API host, temperature, or timeout:

```go
package main

import (
	"context"
	"log"
	"os"
	"time"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
)

func main() {
	ctx := synapseq.NewAppContext()
	loaded, err := ctx.AI(context.Background(), "Generate a 10 minute relaxation sequence", &synapseq.AIOptions{
		Model:       "local-model",
		BaseURL:     "http://localhost:1234",
		Temperature: 1,
		Timeout:     90 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("relaxation.spsq", loaded.RawContent(), 0o644); err != nil {
		log.Fatal(err)
	}
}
```

`AIOptions.Temperature` and `AIOptions.Timeout` belong to the calling
application. The CLI-only environment variables `SYNAPSEQ_AI_TEMPERATURE` and
`SYNAPSEQ_AI_TIMEOUT` are not resolved automatically by the Go API.

Omit `Model` or `BaseURL` to use `SYNAPSEQ_AI_MODEL`,
`SYNAPSEQ_AI_BASE_URL`, and the default OpenAI configuration. Pass a custom
`context.Context` to `AppContext.AI` when the application needs cancellation.
The returned `LoadedContext` is already validated and can also be rendered
with `loaded.WAV`.

## SPSQ Builder

The builder API constructs an SPSQ sequence in code and sends it through the
regular loading, validation, and rendering pipeline:

```go
package main

import (
	"fmt"
	"os"
	"time"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/spsq"
)

func main() {
	ctx := synapseq.NewAppContext().WithVerbose(os.Stderr, true)
	builder, err := spsq.New(ctx)
	if err != nil {
		panic(err)
	}
	builder.
		SampleRate(44100).
		Volume(80).
		Waveform("softpulse", 0, 0, 20, 60, 100, 60, 20, 0)

	focus := builder.NewPreset("focus")
	focus.Tone(220).Waveform("softpulse").Binaural(12).Amplitude(25)
	focus.Pink(15).Amplitude(12)

	timeline := builder.
		SilenceAt(0).
		PresetAt(15*time.Second, focus).
		PresetAt(4*time.Minute+30*time.Second, focus).
		SilenceAt(5 * time.Minute)

	loaded, err := timeline.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(loaded.RawContent()))
	if err := loaded.WAV("output.wav"); err != nil {
		panic(err)
	}
}
```

`Builder.Waveform(name, points...)` emits a top-level `@waveform` definition with 2 through 16384 points in the range `0..100`. `Preset.Waveform(name)` selects that named waveform for the most recently added track. The regular SPSQ loading pipeline validates names, point counts, ranges, duplicates, built-in conflicts, and references when `Load` is called.

`Preset.Shift(separation)` adds `shift` to the most recently added tone, noise, ambiance, or music track. The separation is measured in Hz and `Preset.Intensity(percent)` controls its dry/wet mix.

`Preset.Doppler(rate)` adds `doppler` to the most recently added tone, ambiance, or music track. The rate is measured in Hz; `Preset.Intensity(percent)` controls waveform-shaped pitch or playback-speed movement.

## Convert SBaGen Sequences

The `sbg` package provides the SBaGen conversion flow in Go:

```go
ctx := synapseq.NewAppContext()
converter, err := sbg.New(ctx)
if err != nil {
	panic(err)
}

loaded, err := converter.LoadFile("session.sbg")
if err != nil {
	panic(err)
}

content := loaded.RawContent()
```

Use `converter.LoadContent(sbgContent)` when the SBaGen source text is already
available in memory.

## Packages

- [`core`](https://pkg.go.dev/github.com/synapseq-foundation/synapseq/v4/core) - public application context and loaded sequences.
- [`spsq`](https://pkg.go.dev/github.com/synapseq-foundation/synapseq/v4/spsq) - programmatic SPSQ builder.
- [`sbg`](https://pkg.go.dev/github.com/synapseq-foundation/synapseq/v4/sbg) - SBaGen conversion.
