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

## Example: explain to create

````markdown
## Recommended next task

- Skill: `create-spsq`
- Reason: materialize the explanation as a new valid sequence
- Input files:
  - none
- Suggested output file:
  - `transition-example.spsq`
- Objective:
  - create a practical demonstration of the four transitions.
- Requirements:
  - include `steady`, `ease-in`, `ease-out`, and `smooth`;
  - add didactic comments;
  - validate with `synapseq -test`.
- Constraints:
  - do not modify existing files;
  - use only supported syntax.
- Completion criteria:
  - new file created;
  - validation completed successfully.

## Prompt for the next skill

```text
$create-spsq Create a new transition-example.spsq file that demonstrates
steady, ease-in, ease-out, and smooth in clearly commented phases. Do not use
or modify existing files. Use only supported syntax and validate the new file
with synapseq -test.
```

## Structured context

```yaml
synapseq_task:
  source_skill: explain-spsq
  target_skill: create-spsq
  operation: create_new_sequence
  source_files: []
  output_file: transition-example.spsq
  objective:
    - demonstrate_supported_transitions
  requirements:
    - include_steady
    - include_ease_in
    - include_ease_out
    - include_smooth
    - add_didactic_comments
    - validate_with_synapseq_test
  preserve: []
  changes: []
  constraints:
    - never_modify_existing_files
    - use_supported_syntax_only
  completion_criteria:
    - new_file_created
    - validation_passed
```
````
