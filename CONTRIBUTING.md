# Contributing to SynapSeq

First off, thank you for considering contributing to **SynapSeq**!

This project grows stronger with community support, whether through code, docs, bug reports, build scripts, or new ideas.

## TL;DR (Quick Summary)

The `main` branch contains the latest Go (V4) codebase.  

The legacy C code (V2, inherited from SBaGen) is available in the `v2` branch. If you wish to view or fork the old version, use that branch. **No new contributions will be accepted to the `v2` branch.** It is frozen for historical reference.

---

## Contribution Policy

- Follow Go best practices and maintain clean, readable code.
- Ensure backward compatibility when possible.

## SynapSeq Principles

SynapSeq is guided by a small set of principles that define both its development process and its overall philosophy. Contributions should align with these principles whenever possible:

- **One way to do it** - Prefer a single, clear approach for each task.
- **Clarity over cleverness** - Favor readable, straightforward syntax and implementation over clever or obscure solutions.
- **Intention over syntax** - Focus on what the user is trying to achieve rather than adding unnecessary syntactic complexity.
- **Less options, more focus** - Prioritize essential features and avoid expanding configuration or surface area without a strong reason.

When making changes, use these principles as a filter for API design, user experience, CLI behavior, documentation, and internal implementation choices.

Before making larger code changes, read [ARCHITECTURE](docs/ARCHITECTURE.md). It explains the package boundaries, runtime flow, and the architectural invariants contributors are expected to preserve.

If your change touches the `.spsq` or `.spsc` language, parser behavior, timeline rules, or sequence-building semantics, also read [SYNTAX](docs/SYNTAX.md).

### Git Flow Workflow

SynapSeq follows the **Git Flow** branching model:

**Supporting Branches:**

- `feature/*` - New features
- `bugfix/*` - Bug fixes
- `hotfix/*` - Critical fixes for production
- `release/*` - Release preparation

---

## Contributing Sequence Files (`.spsq`)

If you'd like to share your own `.spsq` sequence files with the community, please contribute them to the [SynapSeq Hub Repository](https://github.com/synapseq-foundation/synapseq-hub).

This process is separate from code contributions to the main SynapSeq project.

---

## Commit Convention

We use the **Conventional Commits** format.  
Examples:

- `feat: add new waveform option`
- `fix: correct parsing bug for noise sequences`
- `docs: update README with usage examples`
- `build: add Makefile for macOS`
- `chore: clean up unused code in parser`

Following this format keeps the commit history clear and enables automated changelog generation in the future.

---

## Running Tests

SynapSeq includes unit and integration tests to ensure code quality and prevent regressions.

### Running All Tests

```bash
make test
```

This command runs all unit and integration tests in the project using Go's testing framework.

### Writing Tests

When contributing code, please:

- Add unit tests for new functions and features
- Update existing tests if you modify behavior
- Ensure all tests pass before submitting your PR
- Follow Go testing conventions (files ending in `_test.go`)
- Use table-driven tests when appropriate for better coverage

Example test locations:

- `internal/audio/*_test.go` - Audio processing tests
- `internal/parser/*_test.go` - Parser and syntax tests
- `internal/sequence/*_test.go` - Sequence loading tests

## Pull Request Guidelines

Before opening a PR, please make sure:

- [ ] You wrote clear, descriptive commit messages (see above).
- [ ] All tests pass successfully (`make test`).
- [ ] You added tests for new features or bug fixes.
- [ ] Your changes are limited to the scope of the PR (no unrelated edits).
- [ ] Code follows Go best practices and conventions.
- [ ] Documentation is updated if your changes affect user-facing features.

---

## Not Sure? Open an Issue First!

If you're not sure whether your contribution fits, **open an issue**.

We'll be happy to discuss your idea before you start coding, saving you time and aligning with the project's roadmap.

---

## License

By contributing to SynapSeq, you agree that your contributions will be licensed under the same license as the project.

Please ensure you have the right to submit your code or content under these terms.

---

Thank you for helping make SynapSeq better! Even small contributions; fixing typos, improving docs, or sharing ideas, help this project grow.
