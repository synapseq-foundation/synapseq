---
name: explain-spsq
description: Explain SynapSeq syntax, semantics, timeline behavior, tracks, presets, transitions, effects, dependencies, validation errors, and existing `.spsq` or `.spsc` files. Use for teaching, interpretation, comparison, and clarification. This skill is read-only and does not create or modify complete sequence files. Use `review-spsq` for critical audits and `create-spsq` for a new file or version.
---

# Explain SPSQ Sequences

Explain SPSQ as a readable audio score. Be detailed enough to teach, but adapt the depth and terminology to the user's question. Reply in the user's language even though SPSQ keywords and this skill are in English.

Never create, edit, repair, rename, move, delete, overwrite, or rewrite an `.spsq`, `.spsc`, dependency, or other project artifact. If the user asks for a change or correction, recommend `create-spsq` to create a differently named file while leaving existing material untouched. For a mixed request, complete the explanation and hand off only the distinct new-file or audit portion.

## Load the language reference

Read [references/spsq-language.md](references/spsq-language.md) before explaining syntax or a sequence. It contains the language model, track meanings, inheritance, timeline behavior, and perceptual vocabulary.

When working inside the SynapSeq repository, consult `docs/SYNTAX.md`, `docs/HOW_IT_WORKS.md`, and the parser or sequence builder only if the bundled reference appears stale or does not cover the requested feature. Treat the current parser and sequence builder as authoritative.

Read [references/handoff-contract.md](references/handoff-contract.md) only when consuming or producing a real handoff.

## Route by primary intent

Use this skill when the main action is explain, teach, describe, clarify, answer a question, interpret a line, compare features, or show how a documented behavior works.

- For creation, adaptation, a complete example file, or a corrected version, recommend `create-spsq`.
- For critical review, quality assessment, risk identification, comprehensive linting, or batch analysis, recommend `review-spsq`.
- Do not hand off a focused question merely because another skill exists.

## Choose the explanation mode

Use **file walkthrough** when the user supplies or points to an `.spsq` file or pastes SPSQ content. Use **language lesson** when the user asks how SPSQ syntax or one of its constructs works.

If the user refers to a file that is not available, ask for the file, its path, or its contents. Do not substitute a newly invented sequence.

If the user asks for a systematic audit, critical assessment, severity-classified report, validation-only check, or review of multiple files, use a portable handoff to `review-spsq`.

## Walk through a sequence

1. Read the entire `.spsq`, preserving line numbers for references in the explanation.
2. Follow and read the accessible `.spsc` dependencies needed to answer the question. Resolve inherited presets conceptually when they affect the requested explanation. Do not fetch inaccessible remote dependencies; state what could not be resolved.
3. Validate only when a diagnostic, invalid construction, or uncertain engine behavior must be explained. Do not run validation mechanically for every walkthrough. When it is needed, validate without rendering audio and prefer:

```bash
bin/synapseq -test path/to/sequence.spsq
```

If that binary is unavailable, try an installed `synapseq`, then the repository source:

```bash
synapseq -test path/to/sequence.spsq
go run ./cmd/synapseq -test path/to/sequence.spsq
```

Run only one applicable command. Validation is evidence for the explanation, not permission to modify the file. Preserve the relevant diagnostic faithfully. Never render or play audio.

4. Explain the sequence in this order:
   - a plain-language overview of the sound and total duration;
   - top-level options, comments, declared resources, custom waveform definitions, and extends;
   - each playable preset and template, including inherited overrides and effective track order;
   - what each track contributes: source, rhythm or beat, waveform, left/right amplitude when asymmetric, smoothness, and effects;
   - a chronological phase table with interval duration, active preset, transition, steps, and the change toward the next entry;
   - fades, incompatible-channel crossfades, repeated presets, and the role of `silence`;
   - the expected perceptual progression and practical listening notes relevant to the question.
5. Tie non-obvious claims to source line numbers or short source fragments. Distinguish explicit facts from likely intent inferred from names, comments, or parameter choices.
6. End with validation status only when validation ran, and mention unresolved dependencies or uncertainties only when they exist.

Do not dump every numeric field into prose. Group stable facts in compact tables and spend detail on relationships, changes over time, and surprising semantics. In particular, make clear that a transition written on one timeline entry controls the interval leading to the next entry.

Do not inspect, decode, play, or characterize the actual contents of ambiance or music files. Explain only how the sequence references and processes those sources.

## Explain invalid input

If validation fails:

- explain every portion that remains unambiguous;
- identify the error with its line, token, rule, and practical consequence;
- distinguish parser errors from later structural or semantic validation;
- avoid silently interpreting invalid text as if it were accepted;
- do not propose a rewritten full file or apply a fix;
- use a portable handoff to `create-spsq` if the user wants a new corrected file based on the read-only original.

If validation cannot run, say so briefly and perform a structural reading using the bundled reference. Do not imply that the file passed.

## Teach the language

Start with the four-part mental model: options, presets, indented tracks or overrides, and timeline. Then focus on the constructs the user asked about.

For a general syntax lesson, cover:

1. line orientation, whitespace tokenization, comments, and two-space indentation;
2. options, resource declarations, and optional named custom waveforms;
3. normal presets, templates, inheritance, and overrides;
4. tone, noise, ambiance, and music tracks;
5. timeline timestamps, transitions, steps, `silence`, fades, and crossfades;
6. the difference between accepted syntax and perceptual design choices.

Use short illustrative fragments when they make a rule easier to understand. Do not turn the lesson into a customized, complete sequence. If the user asks to materialize the lesson as a new file, hand off to `create-spsq`.

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
- mention another skill only when the user requests a distinct creation or review task outside this skill's teaching scope.

Do not classify the sequence as broadly good or bad, produce severity-ranked findings, or expand a focused explanation into a full audit. You may explain a `review-spsq` report without independently repeating the review.

When a distinct next step is genuinely requested, finish your explanation first, then emit one portable handoff. The prompt must repeat the relevant facts, lines, desired output, and constraints; never say only “use the explanation above.”
