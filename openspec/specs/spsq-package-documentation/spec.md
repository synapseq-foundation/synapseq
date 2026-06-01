## Requirements

### Requirement: Document spsq package API
The `spsq` package SHALL provide package-level Go documentation that describes it as the programmatic API for building `.spsq` sequence text in Go.

#### Scenario: Package overview appears in Go documentation
- **WHEN** Go documentation is generated for `github.com/synapseq-foundation/synapseq/v4/spsq`
- **THEN** the package documentation explains that callers use `spsq.New()` and the fluent builder to construct `.spsq` content programmatically

#### Scenario: Documentation shows core integration
- **WHEN** a caller reads the `spsq` package documentation
- **THEN** the documentation shows or explains that generated builder output can be passed to `core.AppContext.LoadContent`

### Requirement: Reference spsq in architecture documentation
The architecture guide SHALL mention `spsq` as a public Go API package for programmatically constructing `.spsq` content while preserving parser, sequence loading, and audio rendering ownership in existing packages.

#### Scenario: Contributor reads package responsibilities
- **WHEN** a contributor reads `docs/ARCHITECTURE.md`
- **THEN** the guide identifies `spsq` as a builder API and clarifies that loading, validation, rendering, and preview behavior remain owned by `core` and internal packages

#### Scenario: Contributor reads dependency flow
- **WHEN** a contributor reads the architecture dependency overview
- **THEN** the guide shows or describes the relationship from `spsq` generated content into the existing `core.LoadContent` flow

### Requirement: Reference spsq in README Go API section
The README SHALL mention the `spsq` package in its Go API documentation and link to the package's pkg.go.dev page.

#### Scenario: User looks for Go APIs
- **WHEN** a user reads the README Go API section
- **THEN** they can discover both the `core` package for loading/rendering sequences and the `spsq` package for programmatic sequence construction
