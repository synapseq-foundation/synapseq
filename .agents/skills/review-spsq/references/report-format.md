# SPSQ Review Report Format

## Contents

- [Finding format](#finding-format)
- [Single-file report](#single-file-report)
- [Batch report](#batch-report)
- [Prioritization](#prioritization)
- [Example report](#example-report)

## Finding format

Give every finding one category and a stable identifier:

```text
[E1] Critical error — file.spsq:12
[W1] Technical warning — file.spsq:20
[A1] Artistic observation — preset focus-high
[S1] Optional suggestion — timeline ending
```

For each finding include:

- evidence: line, command output, effective preset, or measured value;
- consequence: why the author should care;
- recommendation: smallest reasonable next action;
- confidence marker when the conclusion is inferred rather than explicit.

Do not repeat a finding in several sections. Refer to its identifier.

## Single-file report

Use this structure, omitting empty sections:

```markdown
# Review of <file>

## Verdict

Direct summary of validity, major risks, and overall structure.

## Validation result

- Command executed:
- Exit status:
- Result:
- File valid: yes / no / not verified

## Critical errors

## Technical warnings

## Sequence structure

| Preset | Tracks | Beat rates | Amplitudes | Effects |
|---|---:|---|---|---|

## Timeline analysis

| Start | Preset | Approx. duration | Transition | Observation |
|---|---|---:|---|---|

## Sound analysis

## Artistic observations

## Optional suggestions

## Prioritized recommendations

1. Essential corrections.
2. Important improvements.
3. Optional experiments.

## Limitations

- Uninspected dependencies or media.
- Static-analysis boundaries.

No source files were modified.
```

When validation does not run, write `File valid: not verified`; never infer “yes” from manual inspection alone.

When the user asks for applied corrections, add one concise note after recommendations:

> This review is read-only. No file was changed. Use `$create-spsq` to create a new corrected version from this source.

## Batch report

Start with:

```markdown
# SPSQ batch review

## Overall summary

| File | Validation | Critical | Warnings | Artistic | Suggestions |
|---|---|---:|---:|---:|---:|
```

Then add one compact section per file, sorted by path:

```markdown
## path/to/file.spsq

- Verdict:
- Validation command and status:
- Critical findings:
- Main technical risks:
- Highest-priority recommendation:
```

Expand only findings specific to that file. Define a recurring issue once, assign it an identifier such as `[COMMON-W1]`, and reference it from subsequent files.

End with ordered cross-file priorities and confirm that no source was modified.

## Prioritization

Order recommendations:

1. validator failures and broken dependencies;
2. valid but unexpected engine behavior;
3. high-impact structural or sound-design risks;
4. responsible-language corrections;
5. optional artistic refinements.

Do not automatically prioritize minimalism, low track count, or one particular entrainment method.

## Example report

```markdown
# Review of psychedelic-journey.spsq

## Verdict

The sequence is structurally dense and intentionally animated. Validation
passes, but two settings that resemble intended options are commented out,
and the simultaneous stereo motion deserves a controlled listening test.

## Validation result

- Command executed: `synapseq -test psychedelic-journey.spsq`
- Exit status: `0`
- Result: `Sequence is valid.`
- File valid: yes

## Technical warnings

[W1] Lines 2–3 contain `# @samplerate 48000` and `# @volume 65`.
They are comments, so SynapSeq uses its defaults of 44100 Hz and volume 100.
If the comments were intended as active settings, the rendered result will
not match that intention.

[W2] All four active tracks use movement effects in `journey-peak`.
This is valid, but competing pan, modulation, and doppler motion may make the
stereo field feel continuously occupied.

## Sequence structure

| Preset | Tracks | Beat rates | Amplitudes | Effects |
|---|---:|---|---|---|
| journey-base | 3 | binaural 7 Hz | 12, 8, 10 | pan, modulation |
| journey-peak | 4 | binaural 12 Hz, isochronic 6 Hz | 14, 10, 12, 8 | pan, modulation, doppler |

## Timeline analysis

| Start | Preset | Approx. duration | Transition | Observation |
|---|---|---:|---|---|
| 00:00:00 | silence | 00:00:30 | smooth | Gradual entrance |
| 00:00:30 | journey-base | 00:07:30 | smooth | Long controlled development |
| 00:08:00 | journey-peak | 00:04:00 | steady | Densest phase |
| 00:12:00 | journey-base | 00:02:30 | smooth | Returns to the base texture |
| 00:14:30 | journey-base | 00:00:30 | smooth | Explicit final fade |
| 00:15:00 | silence | terminal | — | End marker |

## Artistic observations

[A1] The dense central phase is consistent with the psychedelic identity.
Its complexity is a defining characteristic, not an error.

## Optional suggestions

[S1] Consider testing a short grounding preset with lower amplitudes and one
stable, non-moving layer before the final return. This would preserve the
identity while increasing contrast.

## Prioritized recommendations

1. Decide whether the commented sample-rate and volume lines should be active.
2. Audition the peak at low volume and verify stereo clarity.
3. Optionally prototype a grounding phase as a separate new sequence.

No source files were modified.
```
