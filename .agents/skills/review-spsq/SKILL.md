---
name: review-spsq
description: Audit existing SynapSeq `.spsq` files without modifying them. Use for technical reviews, validation, critique, directory or multi-file audits, and analysis of syntax, builder semantics, dependencies, sound structure, timeline composition, responsible claims, or likely unexpected behavior. Classify findings as critical errors, technical warnings, artistic observations, or optional suggestions. Never save corrections; recommend `$create-spsq` when the user wants a new corrected sequence.
---

# Review SPSQ Sequences

Review SPSQ files as a technical auditor and thoughtful sound-design critic. Stay read-only even when the user asks to fix, improve, optimize, or apply changes. Reply in the user's language.

## Load the review references

Read all three references before reviewing:

- [references/review-checklist.md](references/review-checklist.md) for objective syntax, semantics, dependencies, and timeline checks;
- [references/sound-design-guidelines.md](references/sound-design-guidelines.md) for sound structure, artistic judgment, safety, and optional rendered-audio analysis;
- [references/report-format.md](references/report-format.md) for finding categories and individual or batch reports.

When working inside the SynapSeq repository, also consult `README.md`, `docs/SYNTAX.md`, `docs/HOW_IT_WORKS.md`, and the parser or sequence builder when a bundled reference appears stale or incomplete. Treat the current implementation as authoritative.

## Keep every source read-only

Never create, edit, format, rename, move, delete, or change permissions on an `.spsq`, `.spsc`, ambiance, music, or other project file. Suggestions and short corrected fragments belong only in the report.

If the user asks to apply corrections:

1. complete the read-only review;
2. state that no file was changed;
3. show only the smallest useful suggested fragments;
4. recommend `$create-spsq` to create a new corrected file from the reviewed source.

Do not emit a complete rewritten sequence under this skill.

## Resolve the review scope

- For one or more explicit files, review exactly those files.
- For a directory, review only `.spsq` files directly inside it by default.
- Recurse into subdirectories only when the user explicitly requests recursion.
- In batch mode, sort paths for deterministic output and continue after every individual failure.
- If no matching `.spsq` exists, report that fact without inventing input.

For every target:

1. read the complete file with line numbers;
2. locate each `@extends`, `@ambiance`, and `@music`;
3. read accessible local `.spsc` dependencies completely;
4. verify local resource resolution and note inaccessible remote resources;
5. do not characterize uninspected audio from a filename alone.

## Run objective validation

Prefer the bundled validator because it records each file independently and does not stop a batch after the first error:

```bash
bash scripts/validate.sh path/to/one.spsq path/to/two.spsq
```

The script selects one available validator in this order: repository `bin/synapseq`, installed `synapseq`, then repository source through `go run`. Record the exact command, output, and exit status reported for each file.

If the script is unavailable, use one applicable command per file:

```bash
bin/synapseq -test path/to/sequence.spsq
synapseq -test path/to/sequence.spsq
go run ./cmd/synapseq -test path/to/sequence.spsq
```

Never render or play audio as a substitute for `-test`. Never say that a sequence is valid unless an available validator completed successfully. If validation cannot run, label the result **static review only** and state the limitation.

Treat validator diagnostics as critical errors, but do not stop there. The validator establishes accepted syntax and builder semantics; it does not decide whether a valid composition is coherent, efficient, comfortable, or artistically effective.

## Perform the manual audit

Review all five dimensions:

1. syntax;
2. builder semantics and dependencies;
3. sound structure;
4. temporal composition;
5. responsible positioning.

Use the checklist to find objective failures and likely unexpected behavior. Build the effective track layout of inherited presets before comparing consecutive timeline states. Explain when compatible tracks interpolate, incompatible channels crossfade, and active/off channels use boundary fades.

Calculate the total duration and each timeline interval. Remember that the transition and steps written on one entry govern the interval toward the next entry. A repeated identical preset creates a static hold unless it intentionally marks a later fade or another target follows.

Do not use preference as proof. A dense psychedelic sequence may deliberately use several moving layers. Describe the consequence, classify the risk, and preserve the stated identity rather than prescribing minimalism.

## Classify every finding

Assign exactly one category:

- **Critical error**: invalid syntax or semantics, failed validation, missing required dependency, or another issue that prevents reliable loading or rendering.
- **Technical warning**: valid but likely unexpected, fragile, excessively dense, inefficient, or difficult to control.
- **Artistic observation**: a valid characteristic that shapes the experience without being wrong.
- **Optional suggestion**: a goal-dependent improvement that the author may reasonably reject.

Do not duplicate one finding across categories. Cite file and line numbers when possible. Separate explicit facts from inference.

A commented line resembling an option, such as `# @volume 60`, is not automatically invalid because defaults may apply. Classify it as critical only when the declared intent or a resulting broken dependency makes the missing option essential; otherwise report it as a technical warning and explain the active default.

## Review sound and timeline responsibly

Assess simultaneous tracks, nearby carriers, multiple beat methods, amplitudes, noise density, competing movement effects, waveform character, masking, contrast, abrupt changes, long static phases, entrances, development, landing, and final fade.

Do not:

- add amplitude percentages and call the result dB or acoustic loudness;
- claim clipping without measuring rendered audio;
- infer universal comfort or quality;
- promise focus, sleep, healing, drug-like effects, consciousness changes, or brain synchronization;
- treat frequency-to-mental-state associations as guaranteed facts.

Mention headphones for binaural content. Add brief listening-safety advice only when the content or claims make it relevant.

## Optionally inspect rendered audio

Do this only when the user explicitly requests deeper audio analysis.

1. Validate the SPSQ first.
2. Create a uniquely named temporary directory.
3. Render to an explicit WAV path inside that directory; never use the CLI's default output beside the source.
4. If available, use `ffprobe` for duration, channels, and sample rate, and `ffmpeg` analysis filters for measured peak or loudness indicators.
5. Report the exact commands and distinguish measured results from interpretation.
6. Treat samples at or near full scale as evidence to investigate, not automatic proof of audible clipping.
7. Remove only the temporary directory created for this review.

If rendering or inspection tools are unavailable, retain the static review and state which deeper measurements could not be made.

## Deliver the report

Follow [references/report-format.md](references/report-format.md). Lead with the verdict and validation result, then present findings by severity, concise preset and timeline tables, sound analysis, artistic observations, and prioritized recommendations.

For batch review, add an overall summary and one compact subsection per file. Do not repeat generic explanations; define a recurring issue once and reference it.

Always finish with:

- whether validation ran and passed for each file;
- which dependencies or media could not be inspected;
- confirmation that no source file was modified;
- static-analysis limitations;
- `$create-spsq` guidance only when the user requested an applied correction or new corrected version.
