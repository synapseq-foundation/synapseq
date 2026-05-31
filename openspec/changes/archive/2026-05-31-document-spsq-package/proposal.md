## Why

The new `spsq` package exposes a public Go API for building `.spsq` sequence text programmatically, but it is not yet documented for package consumers or visible in the project-level docs. Documenting it now makes the API discoverable alongside the existing `core` and `external` public packages.

## What Changes

- Add package documentation for `spsq` following the style used by `core` and `external`.
- Include a practical Go example based on the temporary `cmd/builder/main.go` usage pattern.
- Update `docs/ARCHITECTURE.md` to describe `spsq` as the programmatic `.spsq` builder API and show its relationship to `core`.
- Update `README.md` to mention the package and link to its pkg.go.dev documentation.
- Avoid changing runtime behavior of the builder API.

## Capabilities

### New Capabilities
- `spsq-package-documentation`: Documents the public programmatic SPSQ builder API and its project-level references.

### Modified Capabilities

## Impact

- Affected packages: `spsq`.
- Affected documentation: `README.md`, `docs/ARCHITECTURE.md`, and package-level Go documentation.
- Public API impact: documentation-only; no intended function, type, or behavior changes.
