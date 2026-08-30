# SPSQ Explanation Reference

## Contents

- [Mental model](#mental-model)
- [Lines, tokens, and comments](#lines-tokens-and-comments)
- [Options and resources](#options-and-resources)
- [Presets, templates, and inheritance](#presets-templates-and-inheritance)
- [Tracks and perceptual meaning](#tracks-and-perceptual-meaning)
- [Timeline and transitions](#timeline-and-transitions)
- [Validation boundaries](#validation-boundaries)
- [Explanation vocabulary](#explanation-vocabulary)

## Mental model

An `.spsq` file is a text score for an evolving audio session:

1. **Options** configure rendering and declare external resources.
2. **Presets** describe sound states.
3. **Tracks and overrides** define or adjust the channels inside those states.
4. **Timeline entries** place the states in time and describe how one changes toward the next.

Normal content appears in that order. Comments and blank lines may appear anywhere. An `.spsc` file is a reusable preset collection: it accepts options and presets but cannot contain a timeline or another `@extends`.

## Lines, tokens, and comments

SPSQ is line-oriented and tokenized by whitespace. It has no quoted-string syntax, so names and local paths cannot contain spaces.

- Options, presets, and timeline entries are top-level.
- Tracks and overrides require exactly two ASCII spaces.
- `# text` is an ordinary structural comment.
- `## text` is also retained as sequence metadata.

Names begin with an ASCII letter, contain only letters, digits, `_`, or `-`, and have at most 20 characters. `silence` is reserved.

The parser first recognizes individual line shapes. The sequence builder then enforces relationships such as ordering, existence of presets, resource resolution, inheritance, and timeline validity. A line can therefore be syntactically recognizable but structurally invalid.

## Options and resources

Supported top-level forms include:

```spsq
@samplerate 44100
@volume 80
@ambiance rain audio/rain
@music bed audio/meditation
@waveform softpulse 0 0 20 60 100 60 20 0
@transition soft-land 0 2 12 42 78 96 100
@extends presets/base
```

- Sample rate defaults to `44100` and must be a positive integer.
- Volume defaults to `100` and accepts `0` through `100`.
- One-argument `@ambiance` and `@music` forms use the resource name as the path.
- Local resource paths are relative, use `/`, omit extensions, and cannot contain `..`.
- Local ambiance tries `.wav` before `.mp3` and loops.
- Local music tries `.mp3` before `.wav` and does not loop automatically.
- Local `@extends` resolves to an `.spsc` file.
- Remote resources must resolve as WAV or MP3 by extension or MIME type.
- `@waveform` defines one reusable cycle with 2 through 16384 decimal points from `0` through `100`. Names are case-sensitive, unique across extensions, and cannot conflict case-insensitively with a built-in waveform.
- `@transition` defines a reusable interpolation curve with 2 through 256 non-decreasing decimal points from `0` through `100`; its first and last points must be `0` and `100`. Names are case-sensitive, unique across extensions, and cannot conflict case-insensitively with a built-in transition.

Options lock when preset, track, override, or timeline content begins. Resource declarations must precede tracks that reference them.

When explaining a sequence, distinguish the declaration from the media itself. An SPSQ file reveals the resource name, location, amplitude, and effects, but not what an uninspected recording contains.

## Presets, templates, and inheritance

A normal preset is a named playable state:

```spsq
focus
  tone 220 binaural 10 amplitude 15
```

A template holds a reusable track layout but cannot be used directly in the timeline:

```spsq
focus-base as template
  tone 220 binaural 8 amplitude 12
  noise pink smooth 25 amplitude 8
```

A derived preset copies an earlier template and adjusts inherited tracks:

```spsq
focus-high from focus-base
  track 1 binaural 12
  track 2 amplitude +2
```

Inherited presets cannot add new tracks. An unsigned numeric override replaces the template value; a value beginning with `+` or `-` changes it relatively. Waveform overrides use an absolute built-in or declared custom name.

Explain both the written override and the effective result. Track order matters because corresponding channels interpolate when consecutive presets have compatible sources and effects. Incompatible channels crossfade at their boundary.

Normal presets must contain tracks, directly or through inheritance. Preset names are unique. Up to 31 user presets coexist with built-in `silence`, and a direct preset has at most 16 track slots. Current override indexes accept `1` through `15`.

## Tracks and perceptual meaning

### Tone

```spsq
  tone 220 amplitude 15
  tone 220 amplitude left 15 right 10
  tone 220 binaural 10 amplitude 15
  tone 220 monaural 10 amplitude 15
  tone 220 isochronic 10 amplitude 15
  waveform triangle tone 220 binaural 10 amplitude 15
  tone 220 effect shift 10 intensity 25 amplitude 15
  tone 220 binaural 10 effect shift 4 intensity 20 amplitude 15
```

- A pure tone is a generated carrier at the stated frequency.
- `binaural` sends nearby frequencies to opposite stereo channels; the stated beat is their difference and stereo headphones are the meaningful listening setup.
- `monaural` mixes nearby frequencies into both channels, producing a physical amplitude beat.
- `isochronic` gates a tone on and off at the stated rate, creating a pronounced pulse.
- Built-in waveforms are `sine` (default), `square`, `triangle`, and `sawtooth`. A custom name may be used after a valid `@waveform` definition.
- Custom values map `0 -> -1`, `50 -> 0`, and `100 -> +1`. Points are equally spaced and linearly joined around the cycle, including final-to-first interpolation.
- The waveform shapes pure, binaural, and monaural oscillators. Isochronic tracks use it for both carrier and gate. Pan, modulation, and doppler also use it for their motion. Shift still processes the resulting waveform, but its quadrature oscillator is always sine/cosine.
- Compatible transitions morph normally between any custom or built-in pair while preserving phase. Sharp segments may add harmonics or aliasing; avoid claiming an exact subjective effect.

Carrier and beat values must be non-negative. Binaural and monaural values must remain below twice the carrier so their lower component stays positive.

### Noise

```spsq
  noise white amplitude 10
  noise pink smooth 20 amplitude 15
  noise brown effect pan 0.1 intensity 30 amplitude 12
  noise pink effect shift 8 intensity 20 amplitude 12
```

- White noise is brightest and evenly distributed by frequency.
- Pink noise is more balanced perceptually, with less high-frequency emphasis.
- Brown noise emphasizes lower frequencies.
- `smooth` reduces moment-to-moment roughness without changing the noise color; it ranges from `0` through `100`.

### Ambiance and music

```spsq
  ambiance rain amplitude 20
	ambiance rain effect doppler 0.8 intensity 40 amplitude 20
  music bed effect modulation 0.1 intensity 25 amplitude 15
	music bed effect doppler 0.8 intensity 40 amplitude 15
  music bed effect shift 10 intensity 25 amplitude 15
```

These tracks play declared external resources. Ambiance loops; music is finite. Explain their configured role without inferring the recording's contents beyond supplied names or comments.

### Effects and levels

- `pan` moves stereo position.
- `modulation` varies amplitude.
- `doppler` adds pitch motion on tones and waveform-shaped playback-speed motion on ambiance and music. It is not supported for noise.
- `shift` supports tone, noise, ambiance, and music. It derives a mono wet signal, shifts wet left by `+value/2` Hz and wet right by `-value/2` Hz, and uses intensity as dry/wet mix. Pure sine shift can resemble a pair with the declared separation; non-sine harmonics move by fixed Hz offsets. Monaural and isochronic tracks shift their complete signal. A binaural track preserves its original pair only in the dry portion and may produce several wet components. Shifted white noise may be heard mainly as stereo decorrelation. It is spectral stereo motion, not a guaranteed binaural beat.
- `intensity` controls effect strength.
- `amplitude VALUE` controls matching left/right levels. `amplitude left LEFT right RIGHT` controls final channel levels independently; `left` always requires `right`. The per-channel gains apply after effects.

Tone supports `pan`, `modulation`, `doppler`, and `shift`. Noise supports `pan`, `modulation`, and `shift`; ambiance and music additionally support `doppler`. On external PCM, Doppler preserves left/right source audio and varies playback speed by up to plus or minus 5 percent at full intensity; each Doppler track has an independent cursor and resumes its current position when returning to fixed speed. Ambiance loops and music can end slightly earlier or later. `shift` uses fixed quadrature regardless of track waveform selection, has a 31-sample wet latency, and may color complex audio or alias near frequency-band edges. Numeric amplitude values for each channel, intensity, smoothness, and volume values are percentages from `0` through `100`. Multiple tracks combine, so a value is not a guarantee of final perceived loudness.

## Timeline and transitions

Timeline form:

```text
HH:MM:SS PRESET [TRANSITION [STEPS]]
```

The first entry is `00:00:00`; later timestamps strictly increase; at least two entries are required. Hours use `00`–`23`, and minutes and seconds use `00`–`59`.

The most important reading rule is directional:

> The transition and steps written on an entry govern the interval from that entry toward the next entry.

For example:

```spsq
00:00:00 silence smooth
00:00:20 focus
```

This means “move smoothly from silence toward `focus` during the first 20 seconds.” It does not mean that `smooth` begins after `00:00:20`.

Transitions:

- `steady` is linear and is the default.
- `ease-in` starts gently and accelerates.
- `ease-out` starts more quickly and settles gently.
- `smooth` eases at both ends.
- A declared custom transition name follows its evenly spaced points; it applies to compatible changes and silence fades, while automatic incompatible-channel crossfades stay linear.

Steps create an alternating forward/backward trajectory before finally arriving at the next state. There are `2 × steps + 1` legs, each needing at least five seconds, with a hard cap of 12 steps. A custom transition curve is reapplied on every leg. Explain steps as repeated oscillation between endpoint tendencies, not as additional timeline entries.

### Silence and fades

Built-in `silence` needs no declaration. Adjacent `silence -> active` and `active -> silence` boundaries become fade-compatible by retaining the active track shape while taking amplitude from or to zero.

To hold a preset and fade only near the end, the active preset is repeated where the fade starts:

```spsq
00:14:00 relax smooth
00:15:00 silence
```

The `smooth` on the `00:14:00` entry controls the final minute toward silence.

When corresponding channels have incompatible source types, effects, ambiance names, or music names, SynapSeq uses an automatic per-channel crossfade around the boundary. It uses up to 30 seconds on each available side and does not insert new timeline periods. Active-to-off and off-to-active channels receive equivalent boundary fades.

The final timestamp is the sequence duration. Each entry begins a period that lasts until the following timestamp; the final entry marks the terminal state rather than adding an unspecified extra duration.

## Validation boundaries

Common invalid conditions include:

- options after preset or timeline content;
- a track without exactly two-space indentation;
- a missing, duplicate, empty, or unknown preset;
- a template referenced directly by the timeline;
- a new track declared inside an inherited preset;
- timeline content inside `.spsc`;
- a first timestamp other than `00:00:00`;
- timestamps that do not strictly increase;
- fewer than two timeline entries;
- a resource path with an extension, absolute root, backslash, or `..`;
- out-of-range values or too many steps for the interval.

When explaining invalid input, preserve the difference between intended meaning and accepted behavior. State what the line appears intended to mean, why SynapSeq rejects it, and what consequence that has. Leave creation of a corrected copy to `$create-spsq`; the original remains read-only.

## Explanation vocabulary

Use these distinctions consistently:

- **Explicit:** directly declared by syntax, such as duration, carrier frequency, or transition.
- **Resolved:** inherited or loaded through an accessible `.spsc`.
- **Inferred:** likely intent suggested by names, comments, or parameter choices.
- **Unknown:** properties of inaccessible dependencies or uninspected media.

Describe audible and structural behavior without promising focus, relaxation, sleep, treatment, cognitive improvement, or brain-state changes. Such labels express creative intent, not guaranteed outcomes.
