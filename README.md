<h1 align="center">SynapSeq</h1>

<p align="center">
<img src="./assets/synapseq-banner-dark.svg" alt="SynapSeq - Neural Audio Sequencing Engine" />
</p>
<p align="center">
  <a href="https://github.com/synapseq-foundation/synapseq/releases/latest"><img src="https://img.shields.io/github/v/release/synapseq-foundation/synapseq?color=blue&logo=github" alt="Release"></a>
  <a href="COPYING.txt"><img src="https://img.shields.io/badge/license-GPL%20v3%20or%20later-blue.svg?logo=open-source-initiative&logoColor=white" alt="License"></a>
  <a href="https://github.com/synapseq-foundation/synapseq/commits"><img src="https://img.shields.io/github/commit-activity/m/synapseq-foundation/synapseq?color=ff69b4&logo=git" alt="Commit Activity"></a>
  <a href="https://skills.sh/synapseq-foundation/synapseq"><img src="https://skills.sh/b/synapseq-foundation/synapseq" alt="skills.sh"></a>
</p>

<p align="center"><strong>SynapSeq - Text-Driven Audio Sequencer for Brainwave Entrainment</strong></p>

**SynapSeq** turns plain-text sequences into evolving audio. Its small domain-specific language lets you combine tones, binaural, monaural, and isochronic rhythms, noise, music, ambiance, effects, and transitions on a precise timeline.

Sequences are stored as readable `.spsq` files, making them easy to inspect, reproduce, share, and keep under version control. You can render them from the command line or build them programmatically with the Go API.

## Why SynapSeq?

Most audio tools represent a session as a visual project file. SynapSeq approaches it as a written score: presets describe **what plays**, while the timeline describes **how the sound changes**.

This makes SynapSeq useful for:

- sound designers and musicians exploring procedural soundscapes;
- creators building sessions for meditation, relaxation, focus, or sleep routines;
- developers integrating deterministic audio generation into Go applications;
- researchers, students, and audio enthusiasts who need experiments to be documented and repeatable;
- communities that want to exchange compact, human-readable audio recipes instead of large project files.

SynapSeq is best understood as a **creative and experimental audio tool**. It gives you precise control over sound generation, but it does not prescribe what a sequence should be used for or promise a particular effect on the listener.

> [!IMPORTANT]
> SynapSeq is not a medical device and is not intended to diagnose, treat, cure, or prevent any condition. Terms such as *brainwave entrainment*, *focus*, *sleep*, and *relaxation* describe common creative or experimental uses, not guaranteed health or cognitive outcomes. Listen at a comfortable volume and use appropriate care when creating or evaluating sessions.

## What You Can Create

- stereo binaural, monaural, and isochronic tone sequences;
- layered noise, music, and ambient soundscapes;
- gradual or stepped changes in pitch, rhythm, amplitude, and other parameters;
- spatial movement and modulation through effects;
- repeatable sessions rendered as WAV, streamed as PCM, played directly, or converted to MP3;
- reusable presets and sequences extended from other `.spsq` files.

## What It Looks Like

A basic `.spsq` sequence is plain text: define options, declare presets with indented tracks, then place presets on a timeline.

```spsq
# Options
@samplerate 44100
@volume 80

# Presets
focus
  tone 220 binaural 12 amplitude 25
  noise pink smooth 15 amplitude 12

# Timeline
00:00:00 silence
00:00:15 focus
00:04:30 focus 
00:05:00 silence
```

See [SYNTAX](docs/SYNTAX.md) for the complete language reference.

Save the example as `focus.spsq`, then render it:

```bash
synapseq focus.spsq
```

The result is a repeatable audio session generated from the text definition. See [HOW IT WORKS](docs/HOW_IT_WORKS.md) for a perceptual explanation of the tone methods, transitions, and effects.

## Quick Start

The recommended way to install SynapSeq is through the platform package manager.

### Homebrew (macOS & Linux)

Install with [Homebrew](https://brew.sh):

```bash
brew tap synapseq-foundation/synapseq
brew trust synapseq-foundation/synapseq # For homebrew >= 6.x
brew install synapseq
```

### Winget (Windows)

Install with [Winget](https://learn.microsoft.com/en-us/windows/package-manager/winget/):

```powershell
winget update
winget install synapseq
```

After installation, you can run `synapseq -install-file-association` to associate `.spsq` files with SynapSeq and enable additional Explorer context menu actions.

### Manual Downloads

If you prefer to install manually, download the appropriate archive from the latest GitHub release: [4.40.1](https://github.com/synapseq-foundation/synapseq/releases/tag/v4.40.1-foundation).

If you want to build SynapSeq from source, see the [Compilation Guide](docs/COMPILE.md).

### Next Steps

After installation on any platform, read the repository docs in this order:

- [SYNTAX](docs/SYNTAX.md)
- [HOW IT WORKS](docs/HOW_IT_WORKS.md)

## SynapSeq Remote

SynapSeq Remote provides ready-to-use sequences. Sync the local index before
listing, searching, downloading, or generating a remote sequence:

```bash
synapseq -sync
```

List all available sequences:

```bash
synapseq -list
```

Search sequences by a word found in their name, description, or category:

```bash
synapseq -search focus
```

Use the sequence ID shown by `-list` or `-search` to download its `.spsq` file.
This is the recommended option when you want to keep and inspect the sequence
on your machine:

```bash
synapseq -download calm-state
```

The file is saved as `calm-state.spsq` in the current directory. You can
also provide a destination directory:

```bash
synapseq -download calm-state ./sequences
```

To download a remote sequence and generate its audio in one step, use `-get`:

```bash
synapseq -get calm-state
```

By default, the output uses the sequence name and the `.wav` extension. An
output file can be specified explicitly:

```bash
synapseq -get calm-state calm-state.wav
```

You can export to mp3 using the `-mp3` flag:

```bash
synapseq -mp3 -get calm-state calm-state.mp3
```

Or with `.mp3` extension:

```bash
synapseq -get calm-state calm-state.mp3
```

## Programmatic API

The public Go API can construct the same `.spsq` representation in code and pass it through the regular loading, validation, and rendering pipeline:

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
	// Create a new app context with colorized verbose logging
	ctx := synapseq.NewAppContext().WithVerbose(os.Stderr, true)
	// Create a new spsq builder with a sample rate of 44100 Hz and volume of 80%
	builder := spsq.New().SampleRate(44100).Volume(80)

	// Create a new preset for focus mode
	focus := builder.NewPreset("focus")
	// Add tone with 220 Hz, binaural with 12 Hz, and amplitude of 25%
	focus.Tone(220).Binaural(12).Amplitude(25)
	// Add pink noise with 15% of smoothness and amplitude of 12%
	focus.Pink(15).Amplitude(12)

	// Create the timeline
	timeline := builder.
		// Fade in 00:00:00
		SilenceAt(0).
		// Focus preset starts at 00:00:15
		PresetAt(15*time.Second, focus).
		// Focus preset ends at 00:04:30
		PresetAt(4*time.Minute+30*time.Second, focus).
		// Fade out at 00:05:00
		SilenceAt(5 * time.Minute)

	// Load the sequence into memory
	loaded, err := timeline.Load(ctx)
	if err != nil {
		panic(err)
	}

	// Print the spsq sequence
	fmt.Println(string(loaded.RawContent()))

	// Save the sequence as a WAV file
	if err := loaded.WAV("output.wav"); err != nil {
		panic(err)
	}
}
```

Docs:
- [core](https://pkg.go.dev/github.com/synapseq-foundation/synapseq/v4@v4.40.3-foundation/core)
- [spsq](https://pkg.go.dev/github.com/synapseq-foundation/synapseq/v4@v4.40.3-foundation/spsq)

## Contributing

We welcome contributions!

Please read the [CONTRIBUTING](CONTRIBUTING.md) file for guidelines on how to contribute code, bug fixes, and documentation to the project.

## License

SynapSeq is distributed under the GPL v3 or later license. See the [COPYING](COPYING.txt) file for details.

### Third-Party Licenses

All original code in SynapSeq is licensed under the GNU GPL v3 or later, but the following components are included and redistributed under their respective terms:

- **[fatih/color](https://github.com/fatih/color)**  
  License: MIT  
  Used for colorized terminal output.

- **[beep](https://github.com/gopxl/beep)**  
  License: MIT  
  Used for audio encoding/decoding.

- **[golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys)**  
  License: BSD 3-Clause  
  Used for platform-specific system integration.

- **[go-colorable](https://github.com/mattn/go-colorable)**  
  License: MIT  
  Used indirectly for cross-platform ANSI color support.

- **[go-isatty](https://github.com/mattn/go-isatty)**  
  License: MIT  
  Used indirectly for terminal capability detection.

- **[pkg/errors](https://github.com/pkg/errors)**  
  License: BSD 2-Clause  
  Used indirectly for error wrapping and stack trace utilities.

All third-party copyright notices and licenses are preserved in this repository in compliance with their original terms.

## Contact

We'd love to hear from you! Here's how to get in touch:

### Issues (Bug Reports & Feature Requests)

Use [GitHub Issues](https://github.com/synapseq-foundation/synapseq/issues) for:

- Bug reports and technical problems
- Feature requests and enhancement suggestions
- Documentation improvements

### Discussions (Questions & Community)

Use [GitHub Discussions](https://github.com/synapseq-foundation/synapseq/discussions) for:

- General questions and support (e.g., "How do I use `@extends`?")
- Help with your sequences (e.g., "My sequence isn't working, can you help?")
- Sharing your own sequences and presets with the community
- Discussing ideas and best practices
- Showcasing creative use cases

### Quick Guidelines

- **Found a bug?** → Open an Issue
- **Want a new feature?** → Open an Issue
- **Need help or have questions?** → Start a Discussion
- **Want to share your sequences?** → Post in Discussions
- **General feedback or ideas?** → Start a Discussion

## Credits

Check out the [CREDITS](CREDITS.md) to see a list of all contributors and special thanks!
