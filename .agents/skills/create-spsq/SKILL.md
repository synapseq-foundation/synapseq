---
name: create-spsq
description: Create, edit, review, and validate SynapSeq `.spsq` brainwave-entrainment sequence files. Use when an agent needs to translate a listening goal into SynapSeq presets and timelines, modify an existing sequence, diagnose SynapSeq DSL errors, work with `.spsc` preset libraries, or check a sequence with the SynapSeq CLI.
---

# Create SPSQ Sequences

Create syntactically valid, listenable SynapSeq sequences while keeping claims about their effects modest. Reply in the user's language even though the DSL keywords and this skill are in English.

## Load the language reference

Read [references/spsq-language.md](references/spsq-language.md) before writing or editing a sequence. It contains the accepted line forms, validated ranges, timeline semantics, and examples.

When working inside the SynapSeq repository, consult `docs/SYNTAX.md` and the parser only if the bundled reference appears stale or the requested feature is not covered. Treat the current parser and sequence builder as authoritative.

## Gather the intent

Extract these requirements from the request:

- intended experience and total duration;
- preferred method: pure tone, binaural, monaural, isochronic, noise, ambiance, music, or a combination;
- headphone availability when considering binaural beats;
- available local paths or URLs for ambiance and music;
- desired output path and whether an existing file must be preserved or edited.

For a vague request such as “make a focus sequence,” ask only for the missing essentials: duration and whether the user prefers a method or delegates that choice. Ask about headphones only when binaural is a likely choice. Do not block on optional details once those essentials are known.

Treat focus, sleep, relaxation, meditation, and similar terms as creative listening goals. Do not promise medical, therapeutic, cognitive, or sleep outcomes. Do not add a generic disclaimer unless it is relevant to what the user asked.

## Design the sequence

1. Outline a small number of phases that fit the requested duration: usually an entrance, one or more active phases, and an exit.
2. Define options first. Add only options the sequence needs; the defaults are sample rate `44100` and volume `100`.
3. Define presets before the timeline. Keep corresponding track purposes in the same declaration order across presets so compatible channels interpolate instead of crossfading.
4. Use `silence` only as the first and/or final timeline entry to mediate a fade-compatible entrance or exit. Never emit consecutive `silence` entries. To begin a fade-out before the end, repeat the active preset at the fade start and place one `silence` entry at the final timestamp.
5. Use `smooth` for rounded changes, `steady` for linear changes, and `ease-in` or `ease-out` only when their directional behavior is intentional.
6. Add steps only for a deliberate back-and-forth trajectory. Ensure the interval is long enough for every leg.
7. Use a template only when two or more playable presets share a base containing multiple tracks and each variant changes only a small subset of parameters. Define presets directly when there is only one variant, only one simple track, most parameters change, or the source/effect types or channel layout differ.
8. Keep amplitudes conservative unless the user supplied exact values. Avoid stacking many high-amplitude tracks merely to make a sequence sound stronger.

Do not invent ambiance or music assets. Use only paths, URLs, or named resources supplied by the user or already present in the sequence and its `.spsc` dependencies. If a requested external layer has no source, ask for it or omit that layer and say so.

## Create or edit the file

For a new sequence, use the requested path. If no path is given and filesystem tools are available, choose a short descriptive lowercase hyphenated filename ending in `.spsq` in the current working directory. Otherwise return one fenced `spsq` block.

For an edit:

1. Read the entire target file and any local `.spsc` files it extends.
2. Preserve unrelated options, comments, resource declarations, preset names, and formatting.
3. Make the smallest coherent change that satisfies the request.
4. Recheck channel ordering and every timeline reference after changing presets.

Use exactly two ASCII spaces for tracks and overrides. Never use quoted strings: the language tokenizes on whitespace and has no quoting syntax.

## Validate and repair

Validate every file-backed result without rendering audio. From the SynapSeq repository, prefer:

```bash
bin/synapseq -test path/to/sequence.spsq
```

If that binary is unavailable, try an installed `synapseq`, then the repository source:

```bash
synapseq -test path/to/sequence.spsq
go run ./cmd/synapseq -test path/to/sequence.spsq
```

Run only one applicable command. On failure, use the diagnostic's file, line, span, found token, expected token, and hint to correct the file, then validate again. Do not render or play audio as a substitute for `-test`.

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

## Deliver the result

State the created or edited path, whether CLI validation passed, and a concise description of the sound sources and phase progression. Mention headphones for binaural content. Do not restate the whole file when it was already written unless the user asks.
