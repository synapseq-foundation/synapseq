# SynapSeq SPSQ Language Reference

## Contents

- [File model](#file-model)
- [Options and resources](#options-and-resources)
- [Presets and inheritance](#presets-and-inheritance)
- [Track forms](#track-forms)
- [Values and limits](#values-and-limits)
- [Timeline behavior](#timeline-behavior)
- [Perceptual guidance](#perceptual-guidance)
- [Complete examples](#complete-examples)
- [Common failures](#common-failures)

## File model

An `.spsq` file is line-oriented and whitespace-tokenized. It has no quoted-string syntax. Blank lines are ignored.

Write content in this order:

1. top-level options;
2. top-level preset declarations;
3. tracks or overrides indented under their preset;
4. top-level timeline entries.

Options lock after the first preset, track, override, or timeline line. Presets, tracks, and overrides cannot appear after the timeline begins.

Comments may appear anywhere:

- `# text` is ignored structurally;
- `## text` is also retained as sequence metadata.

Names must start with an ASCII letter, contain only letters, digits, `_`, or `-`, and be at most 20 characters. Preset names are normalized to lowercase after loading; keep resource-name spelling consistent because ambiance and music lookups use their declared names. `silence` is reserved.

## Options and resources

Supported top-level forms:

```spsq
@samplerate 44100
@volume 80
@ambiance rain audio/rain
@ambiance ocean
@music bed audio/meditation
@music theme
@extends presets/base
```

- Sample rate must be an integer greater than zero; default `44100`.
- Volume must be an integer from `0` through `100`; default `100`.
- The one-argument ambiance/music form uses the resource name as its path.
- A local `@extends` path resolves to `.spsc`, not `.spsq`.
- `.spsc` files use the same options and preset syntax but cannot contain timeline entries or another `@extends`.

Local paths:

- use `/`, never `\`;
- must be relative and contain no `..` segment;
- omit the file extension;
- cannot contain spaces because the DSL has no quoting.

Local ambiance resolves `.wav` first, then `.mp3`; WAV is preferable for seamless loops. Local music resolves `.mp3` first, then `.wav`; music is finite and does not loop automatically. Remote URLs must resolve to WAV or MP3 by extension or MIME type.

Declare every ambiance or music name before referencing it in a track.

## Presets and inheritance

Supported declarations:

```spsq
focus

focus-base as template

focus-strong from focus-base
  track 1 amplitude +10
```

- A normal preset must contain at least one track.
- Template presets cannot appear in the timeline.
- `from` may reference only an earlier template preset.
- An inherited preset cannot declare new tracks; modify inherited tracks with overrides.
- A template itself cannot contain overrides.
- Preset names must be unique. At most 31 user presets can coexist with built-in `silence`.
- A preset has 16 channel slots. Direct track declarations fill them in order.

Use a template only when all of these are true:

- at least two playable presets reuse the same base;
- the base contains multiple tracks whose declaration order and types stay aligned;
- each derived preset changes only a few supported parameters.

Prefer direct preset declarations when there is only one variant, the preset has one simple track, most values would be overridden, or tracks/effect types need to be added, removed, or replaced. Templates reduce meaningful duplication; they are not the default way to declare every related preset.

Override syntax is:

```text
  track INDEX KIND VALUE
```

The current parser accepts `INDEX` values `1` through `15`. Override kinds are `tone`, `binaural`, `monaural`, `isochronic`, `waveform`, `pan`, `modulation`, `doppler`, `smooth`, `amplitude`, and `intensity`.

Numeric overrides beginning with `+` or `-` are relative to the template value; unsigned values replace it. The override must match the inherited track: for example, `smooth` requires noise, `binaural` requires a binaural track, and `pan` requires an existing pan effect. Waveform values are absolute keywords, not numeric.

## Track forms

Indent every track with exactly two ASCII spaces.

### Tones

```spsq
  tone 220 amplitude 15
  tone 220 binaural 10 amplitude 15
  tone 220 monaural 10 amplitude 15
  tone 220 isochronic 10 amplitude 15
  waveform triangle tone 220 binaural 10 amplitude 15
  tone 220 effect pan 0.2 intensity 50 amplitude 15
  tone 220 binaural 10 effect doppler 0.8 intensity 40 amplitude 15
```

Waveforms are `sine` (default), `square`, `triangle`, and `sawtooth`. Tone effects are `pan`, `modulation`, and `doppler`. When present, tokens must occur in the shown order: optional beat, optional effect, `intensity`, then `amplitude`.

### Noise

```spsq
  noise white amplitude 10
  noise pink smooth 20 amplitude 15
  noise brown effect pan 0.1 intensity 30 amplitude 12
  noise pink smooth 30 effect modulation 0.2 intensity 35 amplitude 12
```

Noise colors are `white`, `pink`, and `brown`. Noise effects are `pan` and `modulation`; `doppler` is not accepted. If both are used, `smooth` precedes `effect`.

### Ambiance and music

```spsq
  ambiance rain amplitude 20
  ambiance rain effect pan 0.1 intensity 40 amplitude 20
  music bed amplitude 15
  music bed effect modulation 0.1 intensity 25 amplitude 15
```

Ambiance and music support `pan` and `modulation`, not `doppler`. Their source name must match an `@ambiance` or `@music` declaration. Although the current parser permits a waveform prefix for these sources, avoid it unless preserving an existing sequence because waveform is primarily meaningful for generated tones.

## Values and limits

- Floats must be ordinary decimal tokens. Scientific notation, `NaN`, and `Inf` are rejected.
- Amplitude and effect intensity: `0` through `100` percent.
- Carrier and beat/resonance: greater than or equal to `0`.
- Binaural and monaural beat values must be less than twice the carrier so the lower component remains positive.
- Noise smoothness: `0` through `100`.
- Effect values: greater than or equal to `0`.
- A direct preset can allocate at most 16 tracks. Keep a template at 15 or fewer if all tracks may need overrides because override index 16 is not accepted by the current parser.

The parser validates ranges, not psychoacoustic suitability. Choose modest amplitudes and parameter motion unless the user specifies exact settings.

## Timeline behavior

Timeline form:

```text
HH:MM:SS PRESET [TRANSITION [STEPS]]
```

Example:

```spsq
00:00:00 silence smooth
00:00:20 focus smooth 1
00:10:00 silence
```

- Use exactly two digits for each field.
- Hours range from `00` to `23`; minutes and seconds range from `00` to `59`.
- The first entry must be `00:00:00`; every later timestamp must strictly increase.
- A valid sequence needs at least two timeline entries.
- Timeline entries must reference existing, non-template presets. Built-in `silence` is always available.
- Transitions are `steady` (default), `ease-in`, `ease-out`, and `smooth`.

The transition and steps on one timeline entry control the interval from that entry toward the next entry. Values interpolate across compatible tracks aligned by channel/declaration order.

Steps create an alternating forward/backward trajectory with `2 * steps + 1` legs before arriving at the next state. Every leg requires at least five seconds, and steps have a hard cap of 12. For an interval of `D` seconds:

```text
max steps = min(12, max(0, floor((floor(D / 5) - 1) / 2)))
```

Examples: 10 seconds permits `0`, 15 seconds permits `1`, 25 seconds permits `2`, and 55 seconds permits `5`.

`silence -> active` and `active -> silence` become fade-compatible boundaries. Incompatible source types, effects, ambiance names, or music names on the same channel use automatic per-channel crossfades of up to 30 seconds on each available side. `active -> off` and `off -> active` use equivalent boundary fades. No timeline entries are inserted.

For generated sessions, use built-in `silence` only at the timeline boundaries:

- put it first when the session should fade in;
- put it last when the session should fade out;
- never write consecutive `silence` entries;
- do not use `silence` as a duration marker or an intermediate bridge between active presets.

Because a transition belongs to the entry where it is written and targets the next entry, start an ending fade by repeating the active preset, then target one final `silence`:

```spsq
00:00:00 silence smooth
00:01:00 relax-deep
00:14:00 relax-deep smooth
00:15:00 silence
```

Here the sequence holds `relax-deep` until `00:14:00`, then fades to silence over the final minute. Writing `00:14:00 silence` followed by `00:15:00 silence` would instead create a redundant silent-to-silent interval; the preceding transition would already have reached silence at `00:14:00`.

## Perceptual guidance

- Binaural uses different left/right frequencies and is most meaningful with headphones.
- Monaural mixes nearby tones into a physically present amplitude beat.
- Isochronic gates one tone on and off for a pronounced pulse.
- White noise is brightest; pink is more balanced; brown emphasizes lower frequencies.
- Higher noise `smooth` values reduce moment-to-moment roughness without changing noise color.
- `pan` moves the stereo position, `modulation` varies amplitude, and `doppler` adds subtle pitch motion.
- `steady` changes uniformly, `ease-in` starts gently, `ease-out` settles gently, and `smooth` eases at both ends.

Use these as descriptive design tools, not medical claims.

## Complete examples

### Minimal sequence

```spsq
@volume 70

focus
  tone 220 binaural 10 amplitude 15
  noise pink smooth 20 amplitude 10

00:00:00 silence smooth
00:00:20 focus
00:04:40 focus smooth
00:05:00 silence
```

The repeated `focus` marks the beginning of the fade-out. The one final `silence` is both the fade target and the sequence end.

### Template and overrides

```spsq
base-focus as template
  tone 220 binaural 8 amplitude 12
  noise pink smooth 25 amplitude 8

focus-low from base-focus
  track 1 binaural 6

focus-high from base-focus
  track 1 binaural 12
  track 2 amplitude +2

00:00:00 silence smooth
00:00:20 focus-low smooth
00:05:00 focus-high smooth
00:09:40 focus-low smooth
00:10:00 silence
```

## Common failures

- A top-level `tone`, `noise`, `ambiance`, `music`, or `track` line: indent it with exactly two spaces under a preset.
- An option below a preset: move every option to the top.
- A path such as `audio/rain.wav`: remove the local extension (`audio/rain`).
- An empty preset: add a track or inherit from a populated template.
- Tracks under a `from` preset: replace them with valid overrides or use a direct preset.
- A template in the timeline: create a non-template preset with `from`.
- A single timeline timestamp: add an end state; the last timestamp defines sequence duration.
- Consecutive final `silence` entries: repeat the active preset at the intended fade start and keep only one final `silence`.
- An intermediate `silence`: transition directly between active presets; SynapSeq handles incompatible channels with automatic boundary crossfades.
- Excessive steps: reduce them or lengthen the following interval.
- Unexpected extra tokens: restore the exact token order; arbitrary inline comments are not supported.
- Spaces in paths: use a path without spaces because quoting does not exist.
