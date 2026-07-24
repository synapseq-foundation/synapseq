---
name: explain-spsq
description: Explain how SynapSeq `.spsq` sequences work and teach the SPSQ language without creating or editing files. Use when an agent needs to walk through a supplied sequence, describe its options, presets, tracks, inheritance, timeline, transitions, steps, fades, resources, or likely listening experience; explain SPSQ syntax or semantics; or identify and explain problems in an invalid sequence. Do not use to create files; when the user wants changes or corrections, recommend `$create-spsq` to create a new sequence from the existing read-only reference.
---

# Explain SPSQ Sequences

Explain SPSQ as a readable audio score. Be detailed enough to teach, but adapt the depth and terminology to the user's question. Reply in the user's language even though SPSQ keywords and this skill are in English.

Never create, edit, repair, or rewrite an `.spsq` or `.spsc` file. If the user asks for a change or correction, recommend `$create-spsq` to create a new file while leaving the existing material untouched. For a mixed request, explain the existing material and redirect only the new-file part.

## Load the language reference

Read [references/spsq-language.md](references/spsq-language.md) before explaining syntax or a sequence. It contains the language model, track meanings, inheritance, timeline behavior, and perceptual vocabulary.

When working inside the SynapSeq repository, consult `docs/SYNTAX.md`, `docs/HOW_IT_WORKS.md`, and the parser or sequence builder only if the bundled reference appears stale or does not cover the requested feature. Treat the current parser and sequence builder as authoritative.

## Choose the explanation mode

Use **file walkthrough** when the user supplies or points to an `.spsq` file or pastes SPSQ content. Use **language lesson** when the user asks how SPSQ syntax or one of its constructs works.

If the user refers to a file that is not available, ask for the file, its path, or its contents. Do not substitute a newly invented sequence.

## Walk through a sequence

1. Read the entire `.spsq`, preserving line numbers for references in the explanation.
2. Follow each accessible `@extends` dependency and read the complete `.spsc`. Resolve inherited presets conceptually so the user can understand the effective tracks. Do not fetch inaccessible remote dependencies; state what could not be resolved.
3. Validate without rendering audio. From the SynapSeq repository, prefer:

```bash
bin/synapseq -test path/to/sequence.spsq
```

If that binary is unavailable, try an installed `synapseq`, then the repository source:

```bash
synapseq -test path/to/sequence.spsq
go run ./cmd/synapseq -test path/to/sequence.spsq
```

Run only one applicable command. Validation is evidence for the explanation, not permission to modify the file. Never render or play audio.

4. Explain the sequence in this order:
   - a plain-language overview of the sound and total duration;
   - top-level options, comments, declared resources, and extends;
   - each playable preset and template, including inherited overrides and effective track order;
   - what each track contributes: source, rhythm or beat, waveform, amplitude, smoothness, and effects;
   - a chronological phase table with interval duration, active preset, transition, steps, and the change toward the next entry;
   - fades, incompatible-channel crossfades, repeated presets, and the role of `silence`;
   - the likely perceptual progression and practical listening notes.
5. Tie non-obvious claims to source line numbers or short source fragments. Distinguish explicit facts from likely intent inferred from names, comments, or parameter choices.
6. End with validation status, unresolved dependencies, or uncertainties only when any exist.

Do not dump every numeric field into prose. Group stable facts in compact tables and spend detail on relationships, changes over time, and surprising semantics. In particular, make clear that a transition written on one timeline entry controls the interval leading to the next entry.

Do not inspect, decode, play, or characterize the actual contents of ambiance or music files. Explain only how the sequence references and processes those sources.

## Explain invalid input

If validation fails:

- explain every portion that remains unambiguous;
- identify the error with its line, token, rule, and practical consequence;
- distinguish parser errors from later structural or semantic validation;
- avoid silently interpreting invalid text as if it were accepted;
- do not propose a rewritten full file or apply a fix;
- recommend `$create-spsq` if the user wants a new corrected file based on the read-only original.

If validation cannot run, say so briefly and perform a structural reading using the bundled reference. Do not imply that the file passed.

## Teach the language

Start with the four-part mental model: options, presets, indented tracks or overrides, and timeline. Then focus on the constructs the user asked about.

For a general syntax lesson, cover:

1. line orientation, whitespace tokenization, comments, and two-space indentation;
2. options and resource declarations;
3. normal presets, templates, inheritance, and overrides;
4. tone, noise, ambiance, and music tracks;
5. timeline timestamps, transitions, steps, `silence`, fades, and crossfades;
6. the difference between accepted syntax and perceptual design choices.

Use short illustrative fragments when they make a rule easier to understand. Do not turn the lesson into a customized, complete sequence. If the user asks to convert the lesson into a new file, recommend `$create-spsq`.

## Keep claims grounded

- Describe binaural, monaural, isochronic, noise, modulation, spatial movement, and transitions as audio behavior.
- Mention that binaural material is most meaningful with stereo headphones.
- Treat focus, relaxation, meditation, sleep, and similar labels as creative intent or likely listening character.
- Do not promise medical, therapeutic, cognitive, sleep, or brainwave outcomes.
- Clearly label interpretations of preset names, comments, and frequencies as inference rather than guaranteed purpose.

## Deliver the explanation

Prefer a layered answer that a beginner can enter easily and an experienced user can inspect:

- lead with what the sequence does overall;
- use one compact timeline table when there are multiple phases;
- explain specialized terms at first use;
- cite relevant lines rather than reproducing the entire file;
- mention `$create-spsq` only when the user requests a new sequence, including a new adapted or corrected version of existing read-only material.
