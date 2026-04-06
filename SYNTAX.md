# SynapSeq Syntax

This document describes the `.spsq` sequence format, how the parser classifies lines, and which semantic rules are enforced by the sequence builder after parsing.

It complements [ARCHITECTURE.md](ARCHITECTURE.md), which focuses on package boundaries, runtime flow, and code-level responsibilities.

## The `.spsq` Format

SynapSeq sequence files are line-oriented text documents. The parser does not implement a general-purpose grammar with nested blocks or quoted string handling. Instead, it classifies each line by leading token and indentation, then the sequence builder enforces placement and semantic rules.

At a practical level, a normal `.spsq` file is divided into four phases:

1. top-level options;
2. preset declarations;
3. indented track declarations or track overrides under presets;
4. top-level timeline entries.

The file can also contain comments and blank lines anywhere.

## Line Classification Order

Each non-empty line is evaluated in this order:

1. comment;
2. option;
3. preset declaration;
4. track declaration;
5. track override declaration;
6. timeline declaration;
7. invalid syntax.

That order matters. For example, track lines are only recognized when they are indented with exactly two leading spaces. Timeline lines are only recognized at top level and only when their first token parses as `HH:MM:SS`.

```mermaid
flowchart TD
    Line[raw line] --> Empty{empty or whitespace?}
    Empty -->|yes| Skip[ignore line]
    Empty -->|no| Comment{starts with # token?}
    Comment -->|yes| CommentHandle[handle comment]
    Comment -->|no| Option{raw line starts with @?}
    Option -->|yes| OptionHandle[parse option]
    Option -->|no| Preset{top-level name?}
    Preset -->|yes| PresetHandle[parse preset]
    Preset -->|no| Track{two-space indented track?}
    Track -->|yes| TrackHandle[parse track]
    Track -->|no| Override{two-space indented track override?}
    Override -->|yes| OverrideHandle[parse override]
    Override -->|no| Timeline{top-level HH:MM:SS?}
    Timeline -->|yes| TimelineHandle[parse timeline]
    Timeline -->|no| Invalid[diagnostic error]
```

## Tokenization Rules

The parser tokenizes by whitespace.

- there is no quoted-string syntax;
- tokens keep their raw textual form until converted by a specific parser rule;
- spans are tracked for each token so diagnostics can point to exact columns.

Numeric parsing is intentionally strict:

- floats reject `NaN` and `Inf`;
- floats reject scientific notation such as `1e10`;
- integer parsing is strict and used in places such as timeline steps and track override indexes.

## Comments

There are two comment behaviors:

- lines beginning with `#` are parser comments and are ignored structurally;
- lines beginning with `##` are also captured as sequence comments and later exposed through `LoadedContext.Comments()`.

That distinction is important because not every comment in a file becomes user-visible metadata.

## Options

Options are top-level lines that begin with `@`.

The currently supported `.spsq` options are:

- `@samplerate <value>`
- `@volume <value>`
- `@ambiance <name> <path-or-url>`
- `@ambiance <name>` as shorthand, where the path defaults to the same name
- `@extends <path-or-url>`

Options must appear before presets or timeline entries. Once the builder has seen a preset, track, override, or timeline line, options are locked.

Local option paths are normalized and validated before use:

- remote URLs are allowed;
- local paths must use `/`, not `\`;
- local paths must be relative;
- absolute paths and Windows drive paths are rejected;
- parent directory traversal such as `..` is rejected;
- local paths must not include file extensions;
- local ambiance paths resolve to `.wav`;
- local extends paths resolve to `.spsc`.

## Presets

A top-level non-indented identifier line is treated as a preset declaration.

Supported forms are:

```text
alpha
base-focus as template
focus-strong from base-focus
```

Rules:

- preset names must pass name validation;
- duplicate preset names are rejected;
- presets must appear before any timeline entries;
- a preset can inherit only from another preset marked `as template`;
- template presets cannot be used directly in the timeline.

The builder inserts a built-in `silence` preset up front, which is why user-defined presets begin after that implicit baseline.

## Track Declarations

Tracks must be indented with exactly two leading spaces under the current preset.

Examples:

```text
alpha
  noise pink amplitude 30
  tone 200 binaural 10 amplitude 15
  waveform square tone 300 isochronic 10 amplitude 8
  ambiance rain amplitude 25
```

Supported track families are:

- `tone`
- `noise`
- `ambiance`

Tone lines can describe:

- pure tones;
- binaural beats;
- monaural beats;
- isochronic beats;
- optional waveform selection via a leading `waveform` token;
- optional effects followed by `intensity` and then `amplitude`.

Noise lines can describe white, pink, or brown noise, optionally with `smooth`, optional effects, and then `amplitude`.

Ambiance lines reference a named ambiance option and then define amplitude, with optional supported effects. The current parser also accepts a leading `waveform` token before ambiance declarations, even though waveform selection is primarily a tone-oriented concept.

Track declarations are rejected when:

- they appear before any preset;
- they appear after timeline entries have started;
- they are attached to a preset that uses `from`, because inherited presets cannot define new tracks.

## Track Overrides

Track overrides also use two-space indentation, but start with the `track` keyword.

Example:

```text
focus-strong from base-focus
  track 1 amplitude 35
```

Overrides are allowed only when all of the following are true:

- the current preset exists;
- the current preset inherits from another preset with `from`;
- the current preset is not itself a template;
- the override appears before timeline entries.

The parser accepts override kinds such as:

- `tone`
- `binaural`
- `monaural`
- `isochronic`
- `waveform`
- `pan`
- `modulation`
- `doppler`
- `smooth`
- `amplitude`
- `intensity`

Numeric overrides may be absolute or relative. Relative overrides are recognized by a leading `+` or `-` in the raw token.

## Timeline Entries

Timeline lines are top-level and start with an `HH:MM:SS` timestamp.

Format:

```text
00:00:00 silence
00:00:20 alpha steady 5
00:02:00 beta smooth 5
```

Supported transition values are:

- `steady`
- `ease-in`
- `ease-out`
- `smooth`

Rules:

- the first timeline entry must start at `00:00:00`;
- time fields must use exactly two digits each;
- hours must be `00` to `23`;
- minutes and seconds must be `00` to `59`;
- timeline entries must reference an existing non-template preset;
- each new timeline entry must be strictly later than the previous one;
- at least two periods must exist in the final sequence.

When a timeline line includes a transition, an optional integer step count may follow. Step counts must be non-negative.

## Structural Placement Rules

The builder applies additional rules after the parser recognizes a line:

- options must stay at the top of the file;
- presets must be declared before any timeline entries;
- tracks and overrides must belong to the most recently declared preset;
- a file with no user presets is invalid;
- an empty preset is invalid;
- a file with fewer than two timeline periods is invalid.

This is why a file can be syntactically parseable line by line and still fail final sequence construction.

## `.spsc` Extends Files

`@extends` targets `.spsc` files, not `.spsq` files.

Those files are parsed through the same line parser but with a different builder mode:

- they may define options and presets;
- they may define tracks and track overrides under presets;
- they may not define timeline entries;
- they may not use `@extends` themselves.

The purpose of `.spsc` is to contribute reusable options, templates, and presets into a main `.spsq` sequence.

## Minimal Mental Model

The easiest way to read a `.spsq` file is:

1. read top-level `@` options;
2. read preset headers;
3. read two-space indented track content under each preset;
4. read top-level time entries that arrange presets into playback periods.

That mental model matches both the parser and the sequence builder.
