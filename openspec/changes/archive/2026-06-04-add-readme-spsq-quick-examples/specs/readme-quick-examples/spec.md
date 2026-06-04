## ADDED Requirements

### Requirement: Show a basic .spsq example in the README
The README SHALL include a compact hand-written `.spsq` example that lets a new user recognize the basic SynapSeq document shape.

#### Scenario: User scans the README for sequence syntax
- **WHEN** a new user reads the README
- **THEN** they can see an example containing at least one preset declaration, one indented track declaration, and timeline entries

#### Scenario: User follows deeper syntax documentation
- **WHEN** a user reads the basic `.spsq` example
- **THEN** the README also points them to the existing syntax documentation for complete language details

### Requirement: Show a basic spsq builder example in the README
The README SHALL include a compact Go example that demonstrates using the public `spsq` builder API to construct and load sequence content.

#### Scenario: Go user scans the README for builder usage
- **WHEN** a Go user reads the README Go API section
- **THEN** they can see an example using `spsq.New()`, `NewPreset`, a track builder call, a timeline call, and `Load(ctx)`

#### Scenario: Go user understands builder responsibility
- **WHEN** a Go user reads the builder example
- **THEN** the example presents the builder as a way to construct `.spsq` content that is loaded through the normal `core.AppContext` path
