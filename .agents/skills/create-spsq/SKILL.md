---
name: create-spsq
description: Create new SynapSeq `.spsq` sequence files from a user description or from an existing sequence used as a read-only reference. Use when the task requires a new file, variant, adaptation, or revised version. Validate every new file, but never modify an existing `.spsq`, `.spsc`, dependency, or occupied destination in place. Use `review-spsq` for critical audits and `explain-spsq` for teaching or clarification.
---

# Create SPSQ Sequences

Create syntactically valid, listenable SynapSeq sequences while keeping claims about their effects modest. Reply in the user's language even though the DSL keywords and this skill are in English.

This is the only SynapSeq skill authorized to create complete `.spsq` files. It creates from scratch or derives a new version from read-only material; it never edits an artifact in place.

## Load the language reference

Read [references/spsq-language.md](references/spsq-language.md) before writing a new sequence. It contains the accepted line forms, validated ranges, timeline semantics, and examples.

When working inside the SynapSeq repository, consult `docs/SYNTAX.md` and the parser only if the bundled reference appears stale or the requested feature is not covered. Treat the current parser and sequence builder as authoritative.

Read [references/handoff-contract.md](references/handoff-contract.md) only when consuming or producing a real handoff to another skill.

## Route by primary intent

Use this skill when the main action is create, generate, produce, assemble, recreate, transform into another file, create a variant, or apply recommendations as a new version.

- For review, audit, quality assessment, risk identification, or validation of an existing file, recommend `review-spsq`.
- For explanation, teaching, line interpretation, or feature comparison, recommend `explain-spsq`.
- For a compound request that requires review findings before creation, hand off to `review-spsq` unless concrete recommendations are already supplied. When recommendations are supplied, create the new file without recreating the audit or extended lesson.

## Gather the intent

Extract these requirements from the request:

- intended experience and total duration;
- preferred method: pure tone, binaural, monaural, isochronic, noise, ambiance, music, or a combination;
- headphone availability when considering binaural beats;
- available local paths or URLs for ambiance and music;
- desired output path;
- any existing `.spsq` supplied as a read-only design or content reference.

When the request contains a `synapseq_task` handoff targeted at `create-spsq`, treat its natural-language prompt as authoritative and the YAML as structured supporting context. Verify the referenced files and requested destination instead of assuming they exist. Preserve the stated objective, constraints, and `preserve` items, but decide the concrete valid SPSQ implementation here.

For a vague request such as “make a focus sequence,” ask only for the missing essentials: duration and whether the user prefers a method or delegates that choice. Ask about headphones only when binaural is a likely choice. Do not block on optional details once those essentials are known.

Treat focus, sleep, relaxation, meditation, and similar terms as creative listening goals. Do not promise medical, therapeutic, cognitive, or sleep outcomes. Do not add a generic disclaimer unless it is relevant to what the user asked.

## Design the sequence

1. Outline a small number of phases that fit the requested duration: usually an entrance, one or more active phases, and an exit.
2. Define options and any explicitly requested custom waveforms first. Add only declarations the sequence needs; the defaults are sample rate `44100` and volume `100`.
3. Define presets before the timeline. Keep corresponding track purposes in the same declaration order across presets so compatible channels interpolate instead of crossfading.
4. Use `silence` only as the first and/or final timeline entry to mediate a fade-compatible entrance or exit. Never emit consecutive `silence` entries. To begin a fade-out before the end, repeat the active preset at the fade start and place one `silence` entry at the final timestamp.
5. Use `smooth` for rounded changes, `steady` for linear changes, and `ease-in` or `ease-out` only when their directional behavior is intentional.
6. Add steps only for a deliberate back-and-forth trajectory. Ensure the interval is long enough for every leg.
7. Use a template only when two or more playable presets share a base containing multiple tracks and each variant changes only a small subset of parameters. Define presets directly when there is only one variant, only one simple track, most parameters change, or the source/effect types or channel layout differ.
8. Keep amplitudes conservative unless the user supplied exact values. When setting distinct left/right amplitudes, use them deliberately as final post-pan channel balance; avoid stacking many high-amplitude tracks merely to make a sequence sound stronger.

Do not invent ambiance or music assets. Use only paths, URLs, or named resources supplied by the user or already present in the sequence and its `.spsc` dependencies. If a requested external layer has no source, ask for it or omit that layer and say so.

## Enforce immutable inputs

An existing artifact means any path present before this task begins. Treat every such path as read-only, including source `.spsq` files, `.spsc` dependencies, ambiance, music, and occupied destinations.

Never overwrite, edit, format, rename, move, delete, change permissions on, or otherwise mutate an existing artifact. Do not use in-place mutation commands such as `sed -i`, `perl -pi`, `truncate`, `rm`, or `mv` on a source, and do not redirect output over an existing path.

If the user asks to edit, change, adapt, extend, or repair an existing `.spsq`, interpret the request as creating a new derived `.spsq`. Read the entire source and every accessible local `.spsc` it extends, but apply the requested changes only to the new output. Preserve unrelated options, comments, resource declarations, preset names, track order, and formatting when copying material into the derivative.

If the user asks only to audit, review, diagnose, or validate an existing file, finish with a portable handoff to `review-spsq` rather than doing a reduced review. If the user asks only for a syntax lesson or didactic walkthrough, hand off to `explain-spsq`.

## Choose a new output path

For a sequence created from scratch, use the requested path only when it does not exist. If no path is given and filesystem tools are available, choose a short descriptive lowercase hyphenated filename ending in `.spsq` in the current working directory.

For a sequence derived from an existing file:

1. Use an explicitly requested output only when it differs from the source and does not already exist.
2. If the requested output is the source or any other existing path, do not write; ask for a new path.
3. If no output is requested, create it beside the source so relative `@extends`, ambiance, and music paths retain their meaning.
4. Prefer `<source-stem>-v2.spsq`; if the source already ends in `-vN`, increment `N`. If that path exists, continue to the next unused version. A descriptive lowercase hyphenated name requested by the user is also acceptable when unused.
5. If the user requests another directory, verify that every relative dependency remains valid from the new location. If SPSQ path rules prevent that, ask for a compatible destination; do not copy or modify dependencies or media.

When filesystem tools are unavailable, return one fenced `spsq` block as a new sequence and state that no file was written.

## Create the file

Recheck that the selected output is still unused immediately before writing. Write only to that reserved path. For a derivative, make the smallest coherent change that satisfies the request and recheck channel ordering and every timeline reference after changing presets.

The newly reserved output may be revised during this creation-and-validation run. That exception applies only to the new output; all artifacts that existed before the task remain immutable.

Use exactly two ASCII spaces for tracks and overrides. Never use quoted strings: the language tokenizes on whitespace and has no quoting syntax.

## Validate the new output

Validate every file-backed result without rendering audio. From the SynapSeq repository, prefer:

```bash
bin/synapseq -test path/to/sequence.spsq
```

If that binary is unavailable, try an installed `synapseq`, then the repository source:

```bash
synapseq -test path/to/sequence.spsq
go run ./cmd/synapseq -test path/to/sequence.spsq
```

Run only one applicable command. On failure, use the diagnostic's file, line, span, found token, expected token, and hint to correct only the new output, then validate again. If the failure belongs to a read-only source or dependency, leave it untouched; adjust the new output without mutating inputs when possible, otherwise report the blocker. Do not render or play audio as a substitute for `-test`.

If validation cannot run because the CLI or referenced media is unavailable, perform the structural checklist below and clearly report that validation was not executed:

- options precede all presets and timeline entries;
- every non-template preset is non-empty;
- tracks and overrides use exactly two spaces;
- inherited presets add no new tracks;
- timeline presets exist and are not templates;
- the first timestamp is `00:00:00`, timestamps strictly increase, and at least two entries exist;
- `silence` appears only at the beginning and/or end, never in consecutive entries;
- a planned fade-out repeats the active preset at the fade start and reaches one final `silence` entry;
- numeric ranges and step limits match the reference;
- local resource paths have no extensions, backslashes, absolute roots, or `..` segments.
- custom waveform definitions use unique non-built-in names, contain 2 through 16384 points in `0..100`, and every reference resolves exactly.

## Deliver the result

State the newly created path, whether CLI validation passed, and a concise description of the sound sources and phase progression. When a source was provided, explicitly confirm that it and its dependencies were used read-only and left unchanged. Mention headphones for binaural content. Do not restate the whole file when it was already written unless the user asks.

For a derived sequence, summarize the meaningful differences from the source without turning the response into a full quality audit.

After completing the new file, recommend `review-spsq` only when an independent audit has clear value, such as for a long sequence, dense layering, many effects, complex inheritance, external dependencies, or a detailed artistic brief. Recommend `explain-spsq` only when a separate didactic walkthrough was requested. For either case, emit the portable handoff contract; never assume another skill can be invoked automatically.
