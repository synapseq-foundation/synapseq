# SPSQ Technical Review Checklist

## Contents

- [Evidence order](#evidence-order)
- [Syntax and placement](#syntax-and-placement)
- [Options and dependencies](#options-and-dependencies)
- [Presets and inheritance](#presets-and-inheritance)
- [Tracks and values](#tracks-and-values)
- [Timeline semantics](#timeline-semantics)
- [Structural quality](#structural-quality)
- [Common valid-but-surprising behavior](#common-valid-but-surprising-behavior)

## Evidence order

Use evidence in this order:

1. current parser and sequence-builder diagnostics;
2. `docs/SYNTAX.md`;
3. bundled review references;
4. comments, names, and apparent author intent.

Never downgrade a validator failure to an artistic preference. Never promote an inference into an objective error.

## Syntax and placement

Check:

- the file is line-oriented and contains no assumed quoted-string syntax;
- options start with `@`, appear at top level, and precede presets or timeline content;
- lines resembling `# @option` are comments, not active options;
- presets and timeline entries are top-level;
- tracks and overrides use exactly two ASCII spaces;
- no track or override appears before a preset or after the timeline begins;
- names begin with an ASCII letter, use only letters, digits, `_`, or `-`, and remain within 20 characters;
- `silence` is not redefined;
- numeric tokens are ordinary decimals, not scientific notation, `NaN`, or `Inf`;
- only documented keywords, effects, waveforms, transitions, and track forms are used;
- token order matches the accepted track form.

Comments beginning with `##` are retained as sequence metadata. Ordinary `#` comments are ignored structurally.

## Options and dependencies

Accepted options are:

- `@samplerate INTEGER`;
- `@volume INTEGER`;
- `@ambiance NAME [PATH-OR-URL]`;
- `@music NAME [PATH-OR-URL]`;
- `@waveform NAME POINT POINT [POINT ...]`;
- `@transition NAME POINT POINT [POINT ...]`;
- `@extends PATH-OR-URL`.

Check:

- sample rate is a positive integer;
- volume is `0` through `100`;
- resource names are valid and unique where required;
- every ambiance/music track references a declaration;
- every declaration that appears intended for playback is actually referenced;
- each custom transition has 2 through 256 non-decreasing `0..100` points, starts at `0`, ends at `100`, and does not use a built-in transition name;
- local paths are relative, use `/`, contain no `..`, spaces, backslashes, roots, drive prefixes, or extensions;
- local ambiance resolves `.wav` before `.mp3`;
- local music resolves `.mp3` before `.wav`;
- local extends resolves to `.spsc`;
- remote resources identify or serve WAV/MP3 where validation can verify them;
- dependencies exist and are readable.
- custom waveform names are valid, unique across the merged sequence, and do not conflict case-insensitively with `sine`, `square`, `triangle`, or `sawtooth`;
- each custom waveform has 2 through 16384 ordinary decimal points, all in `0..100`;
- every custom waveform reference resolves with exact case.

An `.spsc` may contain options and presets but no timeline or nested `@extends`. Do not report a declared-but-unused ambiance or music resource as invalid; classify it as an efficiency warning.

## Presets and inheritance

Check:

- preset names are unique;
- normal presets contain at least one direct or inherited track;
- template presets are used only as inheritance sources;
- a timeline never references a template;
- `from` references an earlier template;
- inherited presets declare overrides, never new tracks;
- templates contain tracks, not overrides;
- override indexes exist in the inherited layout;
- override kinds match the inherited track and effect;
- relative numeric overrides use a leading `+` or `-`;
- waveform overrides use an absolute built-in or declared custom waveform name;
- user presets remain within the engine limit;
- direct presets remain within 16 channels;
- current override indexes remain within `1` through `15`.

Resolve inherited values before reviewing effective amplitude, beat rates, effects, or channel compatibility.

Look for valid but questionable structure:

- templates with no derived presets;
- templates used by only one trivial variant;
- multiple direct presets with substantial identical multi-track layouts;
- presets never used by the timeline or inheritance;
- derived presets whose many overrides obscure rather than reduce duplication;
- redundant presets with identical effective channel states.

These are warnings or optional refactoring suggestions, not syntax errors.

## Tracks and values

Supported source families:

- `tone`;
- `noise`;
- `ambiance`;
- `music`.

Tone may be pure, binaural, monaural, or isochronic. Built-in waveforms are `sine`, `square`, `triangle`, and `sawtooth`; declared custom names are accepted in the same waveform position.

Effects:

- tone: `pan`, `modulation`, `doppler`;
- noise: `pan`, `modulation`;
- ambiance/music: `pan`, `modulation`.

Check:

- carrier and beat/resonance values are non-negative;
- binaural and monaural beats remain below twice the carrier;
- left and optional right amplitude values, and effect intensity, are `0` through `100`;
- a one-value amplitude is interpreted as matching left/right levels; explicit stereo amplitude must use `left VALUE right VALUE`; review asymmetric levels as intentional channel balance after `pan`;
- noise smoothness is `0` through `100`;
- effect values are non-negative;
- `smooth` appears only on noise and before any effect;
- effect is followed by its value, `intensity`, then amplitude;
- external resource names match declarations.
- custom points are interpreted as evenly spaced bipolar values (`0 -> -1`, `50 -> 0`, `100 -> +1`) with circular linear interpolation;
- an isochronic custom waveform intentionally shapes both its carrier and gate, not only the pulse envelope;
- a non-sine waveform on a tone with `doppler` intentionally shapes the pitch movement as well as the source waveform;
- waveform prefixes on ambiance/music affect pan or modulation motion, not the external PCM itself.

Do not invent unsupported codecs, effects, source types, or implied parameters.

## Timeline semantics

Timeline form:

```text
HH:MM:SS PRESET [TRANSITION [STEPS]]
```

Check:

- the first entry is exactly `00:00:00`;
- each field uses two digits and valid hour/minute/second ranges;
- timestamps strictly increase;
- at least two entries exist;
- every referenced preset exists and is playable;
- transitions are `steady`, `ease-in`, `ease-out`, `smooth`, or a declared custom transition;
- steps are non-negative, at most 12, and fit the interval;
- total duration equals the final timestamp.

The transition and steps on an entry control the interval toward the next entry. Steps create `2 × steps + 1` alternating legs, each requiring at least five seconds. A custom curve repeats on each leg; automatic incompatible-channel crossfades remain linear.

Review:

- duration of every interval;
- whether a first `silence` produces an intentional entrance;
- whether an active preset repeated before final `silence` intentionally starts a late fade;
- whether final `silence` provides an intentional landing;
- consecutive or intermediate silence that adds no useful behavior;
- transitions too short to make large parameter changes feel controlled;
- steps that add unnecessary oscillation;
- long holds with no changing parameters;
- abrupt growth without later release;
- comments or phase names contradicted by the actual timeline.

## Structural quality

Compare effective tracks by channel order between every consecutive pair of presets.

Compatible source and effect layouts interpolate. Incompatible source types, effect types, ambiance names, or music names trigger automatic per-channel boundary crossfades of up to 30 seconds on each available side. Active/off boundaries receive equivalent fades.

Check:

- whether track ordering aligns intended counterparts;
- whether an automatic crossfade is expected or accidental;
- whether source changes are masked by short neighboring periods;
- whether the same effective preset repeats with no parameter evolution;
- whether a repeated preset instead serves a valid hold or fade marker;
- whether unused channel slots alternate active/off;
- whether unnecessary preset duplication obscures the narrative.

## Common valid-but-surprising behavior

- Commented options do not apply; defaults may make the file valid anyway.
- Two identical timeline states do not evolve merely because a transition keyword is present.
- A transition belongs to the interval after the line where it is written.
- Music becomes silent when its finite file ends; the SPSQ timeline may continue.
- Ambiance loops and may expose MP3 padding at loop points.
- A source/effect mismatch crossfades instead of interpolating parameters.
- Amplitude is a control percentage, not dB SPL.
- Compatible waveform changes, including custom-to-built-in changes, morph by phase-aligned table interpolation rather than boundary crossfade.
- A syntactically valid combination can still be dense, masked, abrupt, or artistically incoherent.
