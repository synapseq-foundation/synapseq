# SynapSeq AI Compatibility Rules

This document defines **track compatibility constraints** that AI models must follow
when generating SynapSeq `.spsq` sequences.

It complements the following documents:

- ai-basic.md (syntax and structural rules)
- ai-concepts.md (theoretical concepts)

This document focuses on **track transition compatibility between presets**.

This document does not replace the response-format contract from `ai-basic.md`.

When generating an answer for a user:

- always follow the output contract defined in `ai-basic.md`
- return only `.spsq` code in a single fenced code block when the request is in scope
- use the out-of-scope fallback message from `ai-basic.md` when the request is unrelated

Compatibility is evaluated on directly connected timeline pairs.

The transition keyword belongs to the earlier line in that pair and shapes the interval until the next line.

---

## CORE RULE

SynapSeq builds smooth transitions between consecutive timeline presets by interpolating values
of the same track across time.

Because of this, **tracks must remain structurally compatible when one preset transitions directly into another**.

Only parameter values may change during a ramp.

Track **types themselves cannot change**.

If a structural change is required, the transition must pass through `silence` first.

---

## TRACK POSITION CONSISTENCY

Tracks are indexed by their position inside the preset.

Example:

```text
preset-a
  noise pink amplitude 20
  tone 200 binaural 8 amplitude 10
  tone 180 binaural 8 amplitude 8
```

preset-b must preserve the same layout.

Valid example:

```text
preset-b
  noise pink amplitude 25
  tone 200 binaural 6 amplitude 12
  tone 180 binaural 6 amplitude 9
```

Invalid example:

```text
preset-b
  noise pink amplitude 25
  tone 200 binaural 6 amplitude 12
  noise brown amplitude 15
```

Reason:

- track 3 changed from a binaural tone to brown noise
- this is a structural mismatch on the same track position

---

## BEAT TYPE COMPATIBILITY

Tone beat types must remain identical between presets.

Allowed:

binaural → binaural
monaural → monaural
isochronic → isochronic

Not allowed:

binaural → monaural
binaural → isochronic
monaural → binaural
isochronic → binaural

Example of INVALID transition:

```text
preset-a
  tone 200 binaural 8 amplitude 12

preset-b
  tone 200 monaural 8 amplitude 12
```

This changes the beat generation method and is incompatible.

---

## NOISE TYPE COMPATIBILITY

Noise colors must remain the same across presets.

Allowed transitions:

pink → pink
white → white
brown → brown

Not allowed:

pink → brown
pink → white
brown → pink
white → brown

Example of INVALID transition:

```text
preset-a
  noise brown amplitude 15

preset-b
  noise pink amplitude 15
```

This changes the noise color and is incompatible.

---

## EFFECT COMPATIBILITY

Effects applied to a track must remain the same across presets.

Only the **effect parameters** may change.

The **effect type itself cannot change**.

Allowed example:

```text
preset-a
  tone 300 binaural 8 effect modulation 1 intensity 40 amplitude 10

preset-b
  tone 300 binaural 8 effect modulation 2 intensity 60 amplitude 10
```

Here the modulation value and intensity changed, which is valid.

Invalid example:

```text
preset-a
  tone 300 binaural 8 effect doppler 1 intensity 40 amplitude 10

preset-b
  tone 300 binaural 8 effect modulation 0.5 intensity 60 amplitude 10
```

This changes the effect type (doppler → modulation) which is incompatible.

Effects must follow the same rule as track types:
they are structural and must remain consistent.

---

## WAVEFORM COMPATIBILITY

Waveforms may change between presets.

Allowed example:

```text
preset-a
  waveform sine tone 200 binaural 8 amplitude 10

preset-b
  waveform triangle tone 200 binaural 8 amplitude 10
```

This is valid because the **track type remains a tone with the same beat type**.

Waveforms are considered a **shape parameter**, not a structural change.

---

## PARAMETERS THAT MAY CHANGE

The following parameters may change between compatible tracks:

carrier frequency
beat frequency
amplitude
smooth value
effect value
effect intensity
waveform shape

These are considered **continuous parameters** and can be ramped smoothly.

---

## PARAMETERS THAT MUST NOT CHANGE

The following structural properties must remain constant:

track type (tone / noise / ambiance)
beat type (binaural / monaural / isochronic)
noise color (pink / white / brown)
effect type (pan / modulation / doppler)

Changing any of these breaks compatibility.

---

## WHEN STRUCTURAL CHANGES ARE NEEDED

If a sequence requires a structural change, the transition must pass through
the built-in preset `silence`.

Example:

```text
00:05:00 alpha
00:06:00 silence
00:06:30 beta
```

The silence period resets all tracks safely.

---

## AI GENERATION RULES

When generating sequences:

1. Consecutive presets without an intervening `silence` entry must maintain compatible track layouts.
2. Track index N must always represent the same structural track across directly connected presets.
3. Beat types must remain identical across directly connected presets.
4. Noise colors must remain identical across directly connected presets.
5. Effect types must remain identical across directly connected presets.
6. Waveforms may change.
7. Only parameter values should change during direct transitions.

If structural differences are required, use `silence` to reset the system.
