<h1 align="center">SynapSeq</h1>

<p align="center">
  <a href="https://github.com/synapseq-foundation/synapseq/releases/latest"><img src="https://img.shields.io/github/v/release/synapseq-foundation/synapseq?color=blue&logo=github" alt="Release"></a>
  <a href="COPYING.txt"><img src="https://img.shields.io/badge/license-GPL%20v3%20or%20later-blue.svg?logo=open-source-initiative&logoColor=white" alt="License"></a>
  <a href="https://github.com/synapseq-foundation/synapseq/commits"><img src="https://img.shields.io/github/commit-activity/m/synapseq-foundation/synapseq?color=ff69b4&logo=git" alt="Commit Activity"></a>
</p>
<p align="center">
  <a href="https://skills.sh/synapseq-foundation/synapseq/create-spsq"><img src="https://img.shields.io/badge/skills.sh-create--spsq-000000?logo=vercel&logoColor=white" alt="create-spsq on skills.sh"></a>
  <a href="https://skills.sh/synapseq-foundation/synapseq/explain-spsq"><img src="https://img.shields.io/badge/skills.sh-explain--spsq-000000?logo=vercel&logoColor=white" alt="explain-spsq on skills.sh"></a>
  <a href="https://skills.sh/synapseq-foundation/synapseq/review-spsq"><img src="https://img.shields.io/badge/skills.sh-review--spsq-000000?logo=vercel&logoColor=white" alt="review-spsq on skills.sh"></a>
</p>
<p align="center">
<img src="./assets/synapseq-banner-dark.svg" alt="SynapSeq - Neural Audio Sequencing Engine" />
</p>

<p align="center">
  <a href="#quick-start"><strong>Quick Start</strong></a>
  &nbsp;&middot;&nbsp;
  <a href="#convert-sbagen-sequences">SBaGen</a>
  &nbsp;&middot;&nbsp;
  <a href="#synapseq-remote">Remote</a>
  &nbsp;&middot;&nbsp;
    <a href="docs/GO_API.md">Go API</a>
   &nbsp;&middot;&nbsp;
   <a href="docs/AI.md">AI Tools</a>
   &nbsp;&middot;&nbsp;
   <a href="docs/AI.md#ai-agent-skills">AI Agents</a>
  &nbsp;&middot;&nbsp;
  <a href="docs/SYNTAX.md">Syntax</a>
</p>

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
- optional independent left/right track amplitudes;
- spatial movement, modulation, and symmetric frequency shifting through effects;
- reusable custom periodic waveforms defined as amplitude points;
- reusable custom timeline-transition curves defined as interpolation points;
- repeatable sessions rendered as WAV, streamed as PCM, played directly, or converted to MP3;
- reusable presets and sequences extended from other `.spsc` files.

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

### Windows

Install with [Winget](https://learn.microsoft.com/en-us/windows/package-manager/winget):

```powershell
winget update
winget install synapseq
```

After installation, you can run `synapseq -install-file-association` to associate `.spsq` files with SynapSeq.

### macOS

Install with [Homebrew](https://brew.sh):

```bash
brew tap synapseq-foundation/synapseq
brew trust synapseq-foundation/synapseq # For homebrew >= 6.x
brew install synapseq
```

### Linux

Install SynapSeq on Debian, Ubuntu, their derivatives, or Fedora with:

```bash
curl -fsSL https://synapseq.org/install.sh | sudo bash
```

The official installer identifies the supported distribution, configures the package repository, and installs SynapSeq automatically.

### Manual Downloads

If you prefer to install manually, download the appropriate archive from the
[latest GitHub release](https://github.com/synapseq-foundation/synapseq/releases/latest).

If you want to build SynapSeq from source, see the [Compilation Guide](docs/COMPILE.md).

### Next Steps

After installation on any platform, read the repository docs in this order:

- [SYNTAX](docs/SYNTAX.md)
- [HOW IT WORKS](docs/HOW_IT_WORKS.md)
- [AI Tools](docs/AI.md), if you want to generate sequences with AI
- [Go API](docs/GO_API.md), if you want to integrate SynapSeq into a Go application

## Convert SBaGen Sequences

Starting with SynapSeq 4.41.0, you can convert classic [SBaGen](http://uazu.net/sbagen/) `.sbg` sequences into `.spsq` files:

```bash
synapseq -sbg session.sbg
synapseq -sbg session.sbg converted.spsq
synapseq -sbg session.sbg - > converted.spsq
```

When the output is omitted, SynapSeq writes a `.spsq` file alongside the source using the same name. Use `-` to send the converted SPSQ content to standard output.

> [!WARNING]
> SBaGen conversion is experimental. More complex sequences can produce differences or conversion errors because SBaGen and SynapSeq do not have identical sound sources, transition behavior, or media support. Review and validate the generated `.spsq` before using it.

For a Go example using the `sbg` library, see the [Go API Guide](docs/GO_API.md#convert-sbagen-sequences).

## AI Tools

SynapSeq can generate validated SPSQ files through OpenAI-compatible APIs and
provides skills for AI coding agents.

See the [AI Guide](docs/AI.md) for CLI configuration, local models, Codex,
Claude Code, and the `create-spsq`, `explain-spsq`, and `review-spsq` skills.

## SynapSeq Remote

SynapSeq Remote provides ready-to-use sequences. Sync the local index before
listing, searching, downloading, or generating a remote sequence:

```bash
synapseq -sync
```

For self-hosted catalogs and advanced repository configuration, see the
[Remote Guide](docs/REMOTE.md).

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

## Go API

The public Go API can construct, validate, convert, and render the same
`.spsq` representation used by the CLI.

See the [Go API Guide](docs/GO_API.md) for examples covering AI generation,
the SPSQ builder, and SBaGen conversion.

## Contributing

We welcome contributions!

Please read the [CONTRIBUTING](CONTRIBUTING.md) file for guidelines on how to contribute code, bug fixes, and documentation to the project.

## License

SynapSeq is distributed under the GPL v3 or later license. See the [COPYING](COPYING.txt) file for details.

### Third-Party Licenses

See [Third-Party Licenses](docs/THIRD_PARTY_LICENSES.md) for the licenses of
bundled dependencies.

## Contact

We'd love to hear from you! Here's how to get in touch:

Use [GitHub Issues](https://github.com/synapseq-foundation/synapseq/issues) for
bugs, feature requests, and documentation improvements.

Use [GitHub Discussions](https://github.com/synapseq-foundation/synapseq/discussions)
for questions, support, sequence sharing, and ideas.

## Credits

Check out the [CREDITS](CREDITS.md) to see a list of all contributors and special thanks!
