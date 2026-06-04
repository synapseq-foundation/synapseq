## Why

The README currently explains what SynapSeq is and links to the syntax and Go API docs, but it does not show the core `.spsq` shape or the builder API at a glance. Adding compact examples will help new users understand the tool in a few seconds before they follow deeper documentation.

## What Changes

- Add a short `.spsq` sequence example to `README.md` that shows options, a preset, tracks, and timeline entries.
- Add a short Go `spsq.Builder` example to `README.md` that shows equivalent programmatic construction and loading.
- Keep the examples minimal and introductory, with links to existing syntax and Go API documentation for details.
- No runtime behavior, CLI behavior, public API, or `.spsq` syntax changes.

## Capabilities

### New Capabilities
- `readme-quick-examples`: Documents the README requirement for concise first-glance examples of hand-written `.spsq` content and programmatic builder usage.

### Modified Capabilities

## Impact

- Affected files: `README.md`.
- Affected systems: project documentation only.
- No dependency, build, parser, sequence-loading, audio, preview, CLI, or API changes.
