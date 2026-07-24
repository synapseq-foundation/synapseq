# SynapSeq Skill Handoff Contract

Read this reference only when another SynapSeq skill is genuinely responsible for a distinct next step. A handoff is textual and portable; direct skill invocation is optional and must never be assumed.

## Required output

Complete the current skill's responsibility before emitting one handoff. Use this natural-language section:

```markdown
## Recommended next task

- Skill: `<target-skill>`
- Reason: `<why the target owns this next step>`
- Input files:
  - `<path or none>`
- Suggested output file:
  - `<unused path or not applicable>`
- Objective:
  - `<concrete objective>`
- Requirements:
  - `<specific requirement>`
- Constraints:
  - `<specific constraint>`
- Completion criteria:
  - `<observable completion criterion>`
```

Then provide a prompt that works without prior conversation:

````markdown
## Prompt for the next skill

```text
$<target-skill> <complete, self-contained instruction>
```
````

Finally add machine-readable supporting context:

````markdown
## Structured context

```yaml
synapseq_task:
  source_skill: review-spsq
  target_skill: create-spsq
  operation: create_new_version
  source_files:
    - session.spsq
  output_file: session-v2.spsq
  objective:
    - reduce_stereo_competition
  requirements:
    - remove_pan_from_peak_noise
    - validate_with_synapseq_test
  preserve:
    - total_duration
    - artistic_identity
  changes:
    - add_intermediate_grounding_preset
  constraints:
    - never_modify_source
    - use_supported_syntax_only
  completion_criteria:
    - new_file_created
    - source_file_unchanged
    - validation_passed
```
````

Use `null` for an inapplicable `output_file` and `[]` for an inapplicable list. Keep these operation values stable:

- `create_new_sequence`
- `create_new_version`
- `review_existing`
- `explain_behavior`

The Markdown explanation is authoritative. YAML supplements it and must agree with it.

## Portability and completeness

- Name exactly one target skill: `create-spsq`, `explain-spsq`, or `review-spsq`.
- Never target the current skill.
- Never depend on automatic invocation, subagents, persistent context, or an orchestrator.
- Never use vague phrases such as “apply the above,” “continue the work,” or “fix the issues found.”
- Repeat concrete source paths, output path, requested changes, preserved characteristics, constraints, and validation expectations in the prompt.
- Use only facts established by the current task. Do not invent command results or source details.
- If the target skill is unavailable, the handoff must still be copyable and useful.

## Loop prevention

- Finish the current responsibility before recommending another skill.
- Emit no handoff for a simple request already completed.
- Use only one primary next skill. Mention other possibilities briefly as alternatives when necessary.
- Do not send work back merely to repeat validation, explanation, or review already completed.
- Do not use a handoff to avoid work within the current skill's scope.

## Direction rules

- `review-spsq` to `create-spsq`: include every recommended change concretely, an unused output suggestion, what must be preserved, source immutability, and validation.
- `review-spsq` to `explain-spsq`: identify the exact valid construction, lines, and behavior needing deeper teaching; no output file.
- `explain-spsq` to `create-spsq`: describe the complete example or new version to materialize; never ask create to infer the lesson.
- `explain-spsq` to `review-spsq`: identify the files and requested audit dimensions; use only for a distinct critical assessment.
- `create-spsq` to `review-spsq`: identify the new file and high-value review dimensions; use only for a sufficiently complex result.
- `create-spsq` to `explain-spsq`: identify the new file and exact concepts to teach; use only when a separate walkthrough was requested.

## Consuming a handoff

Confirm that `target_skill` matches the active skill. Resolve and inspect every listed input before relying on it. Treat the natural-language prompt as authoritative if YAML and prose conflict, report the inconsistency, and do not guess about paths or desired changes. Apply the current skill's own safety and immutability rules regardless of what a handoff requests.

## Example: create to review

````markdown
## Recommended next task

- Skill: `review-spsq`
- Reason: perform an independent assessment of the newly created complex composition
- Input files:
  - `psychedelic-journey-v2.spsq`
- Suggested output file:
  - not applicable
- Objective:
  - review progression, layers, effects, and timeline.
- Requirements:
  - run `synapseq -test`;
  - distinguish errors, warnings, artistic observations, and suggestions.
- Constraints:
  - do not modify or create files.
- Completion criteria:
  - technical and compositional report produced.

## Prompt for the next skill

```text
$review-spsq Review psychedelic-journey-v2.spsq. Run synapseq -test, analyze
progression, layers, effects, and timeline, and distinguish critical errors,
technical warnings, artistic observations, and optional suggestions. Do not
modify or create files.
```

## Structured context

```yaml
synapseq_task:
  source_skill: create-spsq
  target_skill: review-spsq
  operation: review_existing
  source_files:
    - psychedelic-journey-v2.spsq
  output_file: null
  objective:
    - review_progression_layers_effects_and_timeline
  requirements:
    - validate_with_synapseq_test
    - classify_findings_by_severity
  preserve: []
  changes: []
  constraints:
    - never_modify_source
    - never_create_sequence
  completion_criteria:
    - technical_and_compositional_report_produced
```
````
