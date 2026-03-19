# SynapSeq AI Syntax Guide (Basic)

This document defines the core rules required for an AI system to
generate valid `.spsq` files.

It focuses on the minimal safe subset of the language.

Use this guide when:

- generating new sequences
- explaining basic SynapSeq syntax
- converting other sequence formats into SynapSeq

This guide avoids advanced features such as template presets.

---

# Output Contract

When responding to a user request that is within scope, the AI must output:

- only the `.spsq` code
- inside a single fenced code block
- with no explanation before the code block
- with no explanation after the code block
- with no title, notes, warnings, bullets, or prose

Valid response shape:

```
@volume 75
@samplerate 44100

alpha
  noise pink amplitude 30
  tone 200 binaural 10 amplitude 15

00:00:00 silence
00:00:20 alpha
00:05:00 alpha
00:05:30 silence
```

If the request is outside the scope of SynapSeq sequence generation or brainwave entrainment, the AI must not attempt to answer the unrelated request.

In that case, it must reply with exactly this English sentence and nothing else:

Sorry, it was not possible to fulfill this request. Please try again.

Requests are out of scope when they are not meaningfully about:

- SynapSeq `.spsq` generation
- brainwave entrainment sessions
- meditation, relaxation, focus, sleep, or similar entrainment-oriented audio sequences
- explanation or transformation of SynapSeq sequence syntax

---

# Mental Model

A SynapSeq file has three sections in strict order:

1.  options
2.  presets
3.  timeline

The parser is strict.

Rules:

- options must appear first
- presets must appear before timeline entries
- timeline entries must be last
- track lines belong under presets
- track lines must start with exactly two spaces

---

# Minimal File Structure

Example:

```text
@volume 75
@samplerate 44100

alpha
  noise pink amplitude 30
  tone 200 binaural 10 amplitude 15

00:00:00 silence
00:00:20 alpha
00:05:00 alpha
00:05:30 silence
```

---

# Comments

Two comment styles exist.

Ignored comment:

    # comment

Visible sequence comment:

    ## comment

Comments must start the line.

Inline comments are invalid.

Invalid example:

    @volume 80 # invalid

---

# Options

Options start with `@`.

Supported options:

    @samplerate NUMBER
    @volume NUMBER
    @ambiance NAME PATH_OR_URL
    @extends PATH_OR_URL

Example:

    @samplerate 44100
    @volume 85
    @ambiance rain audio/rain

Local path rules:

- forward slashes only
- no absolute paths
- no `..`
- no file extension

Resolution:

- ambiance → `.wav`
- extends → `.spsc`

---

# Presets

A preset defines a named sound configuration.

Example:

    alpha
      noise pink amplitude 30
      tone 200 binaural 10 amplitude 15

Rules:

- preset names start with a letter
- may contain letters, digits, `_`, `-`
- maximum length: 20 characters
- `silence` is reserved

---

# Track Definitions

Track lines must:

- be under a preset
- use exactly two leading spaces

Example:

    alpha
      tone 300 amplitude 10

One space or three spaces is invalid.

---

# Track Types

## Pure Tone

    tone CARRIER amplitude LEVEL

Example:

      tone 300 amplitude 10

---

## Beat Tone

    tone CARRIER binaural BEAT amplitude LEVEL
    tone CARRIER monaural BEAT amplitude LEVEL
    tone CARRIER isochronic BEAT amplitude LEVEL

Example:

      tone 300 binaural 10 amplitude 15

---

## Noise

    noise white amplitude LEVEL
    noise pink amplitude LEVEL
    noise brown amplitude LEVEL

Optional smoothing:

    noise pink smooth VALUE amplitude LEVEL

Example:

      noise pink smooth 40 amplitude 20

---

## Ambiance

    ambiance NAME amplitude LEVEL

Example:

      ambiance rain amplitude 25

---

# Waveforms

Tone and ambiance tracks may include a waveform prefix.

Supported values:

    sine
    square
    triangle
    sawtooth

Example:

      waveform triangle tone 250 amplitude 10

---

# Effects

Effects appear before amplitude.

General form:

    effect NAME VALUE intensity PERCENT

Example:

      tone 300 binaural 10 effect modulation 6 intensity 40 amplitude 18

Supported effects:

Tone tracks:

    pan
    modulation
    doppler

Noise tracks:

    pan
    modulation

Ambiance tracks:

    pan
    modulation

---

# Timeline

Timeline entries schedule presets.

Format:

    HH:MM:SS PRESET [TRANSITION]

Example:

    00:00:00 silence
    00:00:20 alpha
    00:05:00 beta smooth
    00:10:00 alpha ease-out
    00:12:00 silence

Rules:

- first entry must be `00:00:00`
- times must be strictly increasing
- timestamps cannot repeat
- preset must exist

Allowed transitions:

    steady
    ease-in
    ease-out
    smooth

Default transition is `steady`.

Transition attachment rule:

- the transition keyword on a timeline line applies to the interval from that line until the next timeline line
- it does not describe how the engine arrived from the previous line to the current one
- the next timeline line is the target state reached at the end of the current interval

---

# Transition Semantics

Transitions describe how the engine moves **from the current timeline line
to the next timeline line**.

Therefore:

- the transition keyword belongs to the current period
- the current line defines the starting state for that period
- the next line defines the target state reached at the end of that period
- transitions only have meaning when the next line changes the target preset state

Example:

    00:00:30 alpha_start smooth
    00:04:00 alpha_deep

Meaning:

- the interval from `00:00:30` to `00:04:00` ramps from `alpha_start` to `alpha_deep`
- `smooth` shapes that ramp across the current period

This is different from:

    00:00:30 alpha_start
    00:04:00 alpha_deep smooth

Here, `smooth` belongs to the interval from `00:04:00` to the next timeline line, not to the ramp that ended at `00:04:00`.

Applying a transition when the preset does not change has **no effect**
and should be avoided.

Bad example:

    00:04:00 alpha_deep smooth
    00:09:00 alpha_deep

Reason:

The next line keeps the same preset (`alpha_deep` → `alpha_deep`), so no parameter change occurs during that interval.

Recommended pattern:

    00:04:00 alpha_deep
    00:09:00 alpha_deep ease-out
    00:10:00 silence

Here, `ease-out` is meaningful because it shapes the interval from `00:09:00` to `00:10:00`, where the target changes from `alpha_deep` to `silence`.

Or when transitioning between different presets:

    00:00:30 alpha_start smooth
    00:04:00 alpha_deep

AI models should only emit transition keywords when the next timeline line changes the target preset state.

---

# Silence Preset

`silence` is a built-in preset.

It represents zero audio.

Example:

    00:00:00 silence
    00:00:20 alpha
    00:05:00 silence

Important:

`silence` marks the moment silence is reached.

The fade occurs between the previous entry and the silence timestamp.

Example:

    00:09:00 alpha ease-out
    00:10:00 silence

Meaning:

- the fade-out happens during the interval from `00:09:00` to `00:10:00`
- at `00:10:00`, silence has already been reached

---

# AI Generation Rules

When generating `.spsq` files:

1.  Start with options, then presets, then timeline.
2.  Use exactly two spaces for track lines.
3.  Never place options after presets.
4.  Never place timeline entries before presets.
5.  Use `HH:MM:SS` time format.
6.  Ensure times are strictly increasing.
7.  Do not emit duplicate timestamps.
8.  Prefer regular presets over templates.
9.  Avoid unnecessary complexity.
10. Remember that the transition keyword on a line applies from that line until the next timeline line.
11. Only use transition keywords when the next timeline line changes the target preset state.
12. Do not apply transitions between identical consecutive presets.

---

# Minimal Valid Example

    alpha
      noise pink amplitude 30
      tone 200 binaural 10 amplitude 15

    beta
      noise pink smooth 40 amplitude 20
      tone 180 binaural 8 amplitude 14

    00:00:00 silence
    00:00:20 alpha
    00:05:00 beta smooth
    00:10:00 alpha ease-out
    00:12:00 silence
