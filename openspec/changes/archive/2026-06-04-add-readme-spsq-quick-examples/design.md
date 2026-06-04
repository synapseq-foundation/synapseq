## Context

The README currently introduces SynapSeq, installation, usage links, and the Go API links. A new reader can understand the project category, but must open separate docs to see what `.spsq` looks like or how the public `spsq` builder is used.

This change is documentation-only. The examples must align with the existing parser syntax and the existing public `spsq` builder API, and they should remain short enough to scan quickly.

## Goals / Non-Goals

**Goals:**
- Show a minimal hand-written `.spsq` example near the README introduction or usage area.
- Show a minimal Go builder example near the existing Go API section.
- Keep both examples small, valid, and connected to the existing deeper documentation links.

**Non-Goals:**
- Change `.spsq` syntax, parser behavior, sequence loading, rendering, preview, or CLI behavior.
- Expand the README into a full syntax guide or API tutorial.
- Add new public builder methods or alter the existing `spsq` package API.

## Decisions

- Use compact examples rather than exhaustive examples.
  - Rationale: the user goal is five-second comprehension; detailed usage belongs in `docs/SYNTAX.md` and package docs.
  - Alternative considered: embed a larger tutorial in the README. That would duplicate existing docs and make the README harder to scan.

- Keep the `.spsq` example focused on one preset and a short timeline.
  - Rationale: this shows the document shape: options, preset, indented tracks, timeline.
  - Alternative considered: include ambiance/music file options. That introduces file setup concerns that distract from the basic mental model.

- Keep the builder example focused on `spsq.New()`, `NewPreset`, track construction, timeline construction, and `Load(ctx)`.
  - Rationale: this mirrors the existing public builder API without implying file output or renderer behavior belongs to `spsq`.
  - Alternative considered: show rendering or preview calls in the same snippet. That is useful but shifts the example away from understanding sequence construction.

## Risks / Trade-offs

- Example drift if syntax or builder APIs change -> Keep examples covered by implementation review against existing parser and `spsq` tests.
- README length increases -> Keep snippets short and avoid duplicating full syntax or API docs.
- Generated builder example may imply validation occurs in the builder itself -> Mention or structure the example around `Load(ctx)`, which uses the normal loading pipeline.
