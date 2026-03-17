# SynapSeq Manual

Version 4.0.0

## Name

SynapSeq is a text-driven audio sequencer for brainwave sessions and ambient sound design.

## Synopsis

```text
synapseq [OPTION]... INPUT [OUTPUT]
synapseq -new TYPE
synapseq -preview INPUT
synapseq -play INPUT
synapseq -test INPUT
synapseq -hub-update
synapseq -hub-list
synapseq -hub-search WORD
synapseq -hub-info NAME
synapseq -hub-download NAME [DIR]
synapseq -hub-get NAME [OUTPUT]
synapseq -manual
synapseq -help
synapseq -version
synapseq -manual -no-color
```

For paged reading, use the command that matches your shell environment.

### Linux and macOS

```text
synapseq -manual | less
```

### Windows PowerShell

```text
synapseq -manual | more
```

## Description

SynapSeq is a text-driven audio sequencer for building repeatable brainwave and ambient listening sessions from plain-text score files. Instead of editing waveforms directly, you describe structure, timing, and sound layers in a small domain-specific language and let the renderer generate the final audio.

A typical sequence combines oscillator-based tones, entrainment beats, filtered noise, ambiance layers, waveform selection, and motion effects. These building blocks are arranged into named presets and then placed on a timeline so the session can evolve predictably over minutes or hours.

The parser is strict by design: options belong at the top, presets come next, and timeline entries close the file. This keeps sessions readable, deterministic, and easier to validate before rendering or playback.

SynapSeq can read local files, remote sequence URLs, or standard input. It can render WAV or MP3 output, stream raw PCM to standard output, preview the session as HTML, play directly through ffplay, and fetch published content from the SynapSeq Hub.

For larger libraries, the language supports modular composition through @extends and reusable template presets, making it practical to share preset families across multiple sessions without duplicating track definitions.

## Command Line

The command line accepts one primary action at a time plus supporting flags. When INPUT is present, it may be a local .spsq file, a URL, or - for standard input. When OUTPUT is omitted, SynapSeq derives a default output path from the input name and target format.

The following options are available from the CLI.

### General Options

- **-new TYPE**: Create a starter sequence from a built-in template. Supported types are meditation, focus, sleep, relaxation, and example.
- **-test**: Validate sequence syntax and semantics without generating audio output.
- **-preview**: Render an HTML timeline preview instead of audio.
- **-play**: Render and play the result directly with ffplay.
- **-quiet**: Suppress non-error CLI output.
- **-no-color**: Disable ANSI colors in CLI output. Useful for pipes, logs, and pagers configured without raw-control support.
- **-manual**: Print the full manual.
- **-help**: Show the concise command overview.
- **-version**: Print version, build, and platform information.

### Hub Options

- **-hub-update**: Refresh the local index of sequences available in the SynapSeq Hub.
- **-hub-clean**: Remove cached Hub data from the local machine.
- **-hub-list**: List sequences currently available from the local Hub index.
- **-hub-search WORD**: Search the Hub index for matching sequence names or metadata.
- **-hub-info NAME**: Show metadata and descriptive details for a Hub sequence.
- **-hub-download NAME [DIR]**: Download a Hub sequence and its dependencies. If DIR is provided, files are stored there.
- **-hub-get NAME [OUTPUT]**: Download a Hub sequence and render it in one step. If OUTPUT is provided, it is used as the target file name.

### External Tool Options

- **-ffmpeg-path PATH**: Use a specific ffmpeg executable when exporting MP3.
- **-ffplay-path PATH**: Use a specific ffplay executable when playing audio directly.

### Windows Options

- **-install-file-association**: Associate .spsq files with SynapSeq on Windows.
- **-uninstall-file-association**: Remove the SynapSeq .spsq file association on Windows.

## File Layout

- **Options**: Start with lines beginning with @. They define global settings such as samplerate, volume, ambiance files, and modular imports.
- **Presets**: Define named sound setups. A preset line has no indentation. Tracks under that preset must start with exactly two spaces.
- **Timeline**: Schedule presets over time using HH:MM:SS. Timeline entries are always top-level lines with no indentation.

The first timeline entry must start at 00:00:00, times must be strictly increasing, and a valid sequence needs at least two periods.

```text
@volume 90
@samplerate 44100

focus
  noise white amplitude 10
  tone 240 binaural 16 amplitude 15

00:00:00 silence
00:00:20 focus
00:19:30 focus
00:20:00 silence
```

## Comments

- **# comment**: Ignored by the parser. Useful for notes and documentation.
- **## comment**: Stored as a sequence comment and shown by the CLI before rendering when output is not quiet.
- **Standalone lines only**: Comments must occupy their own lines. Inline comments after options, preset names, track lines, overrides, or timeline entries are invalid syntax.

```text
@samplerate 48000 # samplerate
```

The example above is invalid. Write comments only on otherwise empty lines.

## Sequence Options

- **@samplerate NUMBER**: Set the output sample rate. Default is 44100.
- **@volume NUMBER**: Set the global output volume from 0 to 100. Default is 100.
- **@ambiance NAME PATH_OR_URL**: Register an ambiance source that can be referenced later by track definitions.
- **@extends PATH_OR_URL**: Load presets and options from a modular .spsc file before the main sequence is built.

For local @ambiance and @extends paths, use relative paths with forward slashes and omit the file extension. Local ambiance paths resolve to .wav and local extends paths resolve to .spsc.

Absolute paths, parent traversal, Windows drive prefixes, and stdin are rejected for local modular paths.

```text
@volume 85
@samplerate 48000
@ambiance rain audio/rain
@ambiance ocean https://example.com/ocean.wav
@extends library/common
```

## Preset Names

Named references should be short and predictable. A preset or ambiance name must start with a letter, may contain letters, digits, underscores, or dashes, and may not exceed 20 characters.

Use lowercase names consistently. Presets are normalized to lowercase internally.

## Presets

- **Regular preset**: A plain top-level name creates a new preset.
- **Template preset**: Use as template to define a reusable base that cannot appear directly in the timeline.
- **Inherited preset**: Use from TEMPLATE_NAME to clone a template preset and then override selected track fields.

```text
focus-base as template
  noise white amplitude 15
  tone 240 binaural 14 amplitude 12

focus-deep from focus-base
  track 1 amplitude 25
  track 2 binaural 18
```

A preset that inherits from a template cannot define new track lines. It must use track overrides instead.

## Track Definitions

Track lines belong under a preset and must start with exactly two spaces. SynapSeq allocates tracks automatically in declaration order.

- **Pure tone**: `tone CARRIER amplitude LEVEL`
- **Beat tone**: `tone CARRIER binaural|monaural|isochronic BEAT amplitude LEVEL`
- **Noise**: `noise white|pink|brown [smooth VALUE] amplitude LEVEL`
- **Ambiance**: `ambiance NAME amplitude LEVEL`

Amplitude is required on every concrete track definition. Tone carrier and beat values are positive numbers. Noise smooth ranges from 0 to 100.

```text
alpha
  tone 300 amplitude 10
  tone 220 binaural 8 amplitude 14
  noise pink smooth 35 amplitude 20
  ambiance rain amplitude 25
```

## Sound Concepts

- **tone**: A steady audible carrier frequency with no beat component. Use it when you want a plain oscillator without binaural, monaural, or isochronic pulsing.
- **binaural**: A beat created by feeding slightly different frequencies to the left and right channels. In SynapSeq syntax, the number after binaural is the beat frequency, not the carrier.
- **monaural**: A beat produced by combining close frequencies into a single perceived pulse before it reaches the ears. It still uses a carrier plus a beat value, but behaves differently from binaural playback.
- **isochronic**: A rhythmic single-tone pulse where the sound is sharply gated on and off. Use it when you want a clearly defined pulse rather than two interacting tones.
- **noise**: A broadband texture instead of a pitched oscillator. White, pink, and brown noise differ in spectral balance, and smooth controls how gently the noise evolves over time.
- **ambiance**: A named WAV file or remote WAV resource mixed into the preset. Ambiance is useful for rain, drones, environmental beds, and other non-synth layers.

A practical rule: tone defines the audible pitch, binaural or other beat keywords define how that pitch is modulated or split, amplitude controls loudness, and optional effects add movement.

## Waveforms

You can prefix tone and ambiance tracks with waveform to change the oscillator shape. Valid values are sine, square, triangle, and sawtooth.

```text
shape-demo
  waveform triangle tone 250 monaural 8 amplitude 10
  waveform square ambiance rain amplitude 20
```

## Effects

Effects are optional and appear before amplitude. All effects require a value and an intensity percentage.

- **pan**: Moves the sound across the stereo field over time. The effect value controls the pan rate and intensity controls how far the movement is pushed.
- **modulation**: Applies cyclic amplitude movement. The effect value acts as the modulation rate and intensity controls how deep the volume swing becomes.
- **doppler**: Applies motion-based pitch and stereo movement for tone tracks. It is available only on tone-based tracks and is useful for a drifting or orbiting sensation.
- **Tone effects**: pan, modulation, and doppler are valid for tone-based tracks.
- **Noise effects**: pan and modulation are valid for noise tracks.
- **Ambiance effects**: pan and modulation are valid for ambiance tracks.

```text
fx-demo
  tone 300 binaural 10 effect modulation 6 intensity 40 amplitude 18
  tone 220 effect doppler 0.9 intensity 70 amplitude 14
  noise white effect pan 0.5 intensity 60 amplitude 12
  ambiance rain effect modulation 4 intensity 35 amplitude 20
```

## Track Overrides

Track overrides are only valid inside presets created with from TEMPLATE_NAME. The syntax is 1-based: track 1 targets the first inherited track.

- **tone VALUE**: Change the carrier frequency of an inherited tone track.
- **binaural|monaural|isochronic VALUE**: Change the beat value of an inherited matching beat track.
- **pan|modulation|doppler VALUE**: Change the effect value when that effect already exists on the inherited track.
- **intensity VALUE**: Change effect intensity.
- **amplitude VALUE**: Change the track amplitude.

```text
focus-strong from focus-base
  track 1 amplitude 30
  track 2 tone 250
  track 2 binaural 18
  track 2 intensity 75
```

## Timeline

- **Format**: `HH:MM:SS PRESET_NAME [TRANSITION]`
- **Default transition**: `steady`
- **Other transitions**: `ease-in, ease-out, smooth`

Timeline entries must reference non-template presets. The parser adjusts each period against the next one to build transitions and track interpolation.

- **steady**: Hold the current preset values without applying an easing curve. This is the neutral transition and the default when you omit the transition keyword.
- **ease-in**: Start the transition gently and accelerate into the target state. Useful when you want the next section to arrive progressively rather than immediately.
- **ease-out**: Move quickly at first and then settle more gently as the next state is reached. Useful for soft landings near the end of a section.
- **smooth**: Apply a balanced easing curve across the whole transition. This is usually the best choice for gradual meditative ramps.
- **Compatibility rule**: Consecutive timeline entries cannot reuse the same channel with an incompatible track type, effect type, or ambiance source.
- **Silence bridge**: If you need to switch a channel from one incompatible sound design to another, insert a silence preset between those timeline entries.
- **Direct on/off changes**: A channel should not jump directly between an active track and off; use silence as the bridge state.

```text
00:00:00 silence
00:00:20 focus-light
00:04:00 focus-light smooth
00:07:00 focus
00:19:30 focus
00:20:00 silence

# Good: incompatible presets are separated by silence
00:00:00 silence
00:00:15 doppler-preset
00:00:30 silence
00:00:45 pan-preset
```

## Extended Files

Files loaded by @extends are modular .spsc files. They may contain options, presets, tracks, and track overrides. They may not contain timeline entries, and they may not contain another @extends option.

- **Library file example**: A `.spsc` file is useful for shared preset libraries. Keep reusable templates or presets there, then import them from your main `.spsq` session.

```text
# library/common.spsc
focus-template as template
  tone 240 binaural 10 amplitude 12
  noise pink smooth 25 amplitude 8
```

- **Main sequence example**: Your main `.spsq` file can import that library with `@extends` and then build concrete presets or timeline entries from it.

```text
@extends library/common

focus-light from focus-template
  track 1 amplitude 10
  track 2 amplitude 6

00:00:00 silence
00:00:20 focus-light
00:10:00 silence
```

## Examples

Start with a minimal session first. Once that structure is clear, add more layers or modular reuse.

### Basic Session

A minimal valid file: one preset, one audible section, and silence at the beginning and end.

```text
@volume 90

focus
  tone 240 binaural 10 amplitude 15

00:00:00 silence
00:00:20 focus
00:00:40 silence
```

### Layered Preset

Add noise or ambiance when you want more texture without changing the overall preset structure.

```text
deep-rest
  noise brown smooth 45 amplitude 12
  tone 180 binaural 6 effect modulation 4 intensity 30 amplitude 16
  ambiance rain effect pan 0.3 intensity 35 amplitude 22

00:00:00 silence
00:00:20 deep-rest
00:10:00 silence
```

### Reusable Templates

Use @extends and template inheritance when several sessions share the same preset family.

```text
@extends library/focus-base
@ambiance rain audio/rain

focus-light from focus-template
  track 1 amplitude 8
  track 2 binaural 12

focus-deep from focus-template
  track 1 amplitude 18
  track 2 binaural 18
  track 2 intensity 65
```

## Notes

- **Use silence**: A built-in preset named silence is always available and is ideal for the first or last timeline entry.
- **Indentation matters**: Preset children must use two leading spaces. A top-level tone/noise/ambiance line is invalid syntax.
- **Keep options first**: Once presets or timeline entries start, additional options are rejected.
- **Prefer templates for families**: If several presets share the same structure, declare a template and override only what changes.
- **Track positions are fixed**: Direct transitions compare tracks by position. Track 1 must stay the same kind of sound across directly connected presets, and the same rule applies to every later track slot.
- **Structural changes need silence**: Do not switch a track directly between tone, noise, or ambiance types. If the structure must change, insert the built-in silence preset between those timeline entries.
- **Beat, noise, and effect types must match**: Across directly connected presets, keep binaural versus monaural versus isochronic mode unchanged, keep white versus pink versus brown noise unchanged, and keep pan versus modulation versus doppler unchanged.
- **Waveforms may change**: Waveform shape is treated as a parameter, not a structural type. It may change between otherwise compatible tone or ambiance tracks.

## See Also

Use synapseq -help for a concise command overview, or visit the online documentation for installation guides and broader examples.

```text
synapseq -help
```
