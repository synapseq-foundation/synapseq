## 1. Package Documentation

- [x] 1.1 Add `spsq/doc.go` with package-level documentation following the style of `core/doc.go` and `external/doc.go`.
- [x] 1.2 Include a concise builder example that constructs sequence text and passes `Builder.String()` to `core.AppContext.LoadContent`.

## 2. Project Documentation

- [x] 2.1 Update `docs/ARCHITECTURE.md` to list `spsq` as a public programmatic `.spsq` builder API and clarify its relationship to `core`.
- [x] 2.2 Update the architecture flow or dependency overview to show generated `spsq` content entering the existing `core.LoadContent` pipeline.
- [x] 2.3 Update `README.md` Go API section with a pkg.go.dev link for `github.com/synapseq-foundation/synapseq/v4/spsq`.

## 3. Verification

- [x] 3.1 Run Go formatting on any added Go documentation file.
- [x] 3.2 Run the relevant Go documentation/test check, at minimum `go test ./spsq ./core`, to confirm examples and package docs compile.
