## 1. Doctor Package Implementation

- [x] 1.1 Create `cmd/synapseq/doctor.go` with tool checking logic
- [x] 1.2 Implement `ToolCheck` struct with Name, Installed, Error fields
- [x] 1.3 Add `CheckTool(name string)` function using `exec.LookPath`
- [x] 1.4 Add platform detection for suggested fixes (darwin/linux → brew, windows → winget)
- [x] 1.5 Add `FormatDoctorOutput(checks []ToolCheck)` function using cli styling

## 2. CLI Integration

- [x] 2.1 Add `-doctor` flag in `internal/cli/flags.go`
- [x] 2.2 Handle `-doctor` flag in `cmd/synapseq/dispatch.go` (early exit if set)
- [x] 2.3 Call doctor check and display results

## 3. Testing

- [x] 3.1 Add unit tests for `cmd/synapseq/doctor.go`
- [x] 3.2 Test tool detection with mock implementations
- [x] 3.3 Test platform-specific output formatting