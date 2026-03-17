<h1 align="center">SynapSeq</h1>

<p align="center">
  <a href="https://synapseq.org/">Home</a> |
  <a href="https://hub.synapseq.org/">Examples</a> |
  <a href="https://synapseq.org/docs">Documentation</a>
</p>

<p align="center">
  <p align="center">
  <a href="https://github.com/synapseq-foundation/synapseq/releases/latest"><img src="https://img.shields.io/github/v/release/synapseq-foundation/synapseq?color=blue&logo=github" alt="Release"></a>
  <a href="COPYING.txt"><img src="https://img.shields.io/badge/license-GPL%20v2-blue.svg?logo=open-source-initiative&logoColor=white" alt="License"></a>
  <a href="https://github.com/synapseq-foundation/synapseq/commits"><img src="https://img.shields.io/github/commit-activity/m/synapseq-foundation/synapseq?color=ff69b4&logo=git" alt="Commit Activity"></a>
</p>
</p>

<p align="center"><strong>Synapse-Sequenced Brainwave Generator</strong></p>

SynapSeq is a lightweight engine that sequences audio tones to guide brainwave states like relaxation, focus, and meditation using a simple text-based format.

🌐 **Visit [synapseq.org](https://synapseq.org/) for installation instructions, documentation, FAQ, and more!**

For a local command reference and language overview in Markdown, see [USAGE.md](USAGE.md). You can also open the terminal manual with `synapseq -manual`.

## AI Specification Pack

If you are using SynapSeq with LLMs or prompt pipelines, see [AI](ai/README.md) for the AI specification pack.

It includes:

- [ai/ai-basic.md](ai/ai-basic.md) for output contract and safe syntax generation
- [ai/ai-compatibility.md](ai/ai-compatibility.md) for transition compatibility rules
- [ai/ai-concepts.md](ai/ai-concepts.md) for entrainment and session design guidance

## Contributing

We welcome contributions!

Please read the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on how to contribute code, bug fixes, and documentation to the project.

## License

SynapSeq is distributed under the GPL v2 license. See the [COPYING.txt](COPYING.txt) file for details.

### Third-Party Licenses

SynapSeq makes use of third-party libraries, which remain under their own licenses.  
All original code in SynapSeq is licensed under the GNU GPL v2, but the following components are included and redistributed under their respective terms:

- **[beep](https://github.com/gopxl/beep)**  
  License: MIT  
  Used for audio encoding/decoding.

- **[go-yaml](https://github.com/goccy/go-yaml)**  
  License: MIT  
  Used for YAML parsing and processing.

- **[pkg/errors](https://github.com/pkg/errors)**  
  License: BSD 2-Clause  
  Used indirectly via `beep` for error wrapping and stack trace utilities.

- **[google/uuid](https://github.com/google/uuid)**  
  License: BSD 3-Clause  
  Copyright © 2009-2014 Google Inc.  
  Used for UUID generation and unique identifier handling.

- **[golang.org/x/sys/windows/registry](https://pkg.go.dev/golang.org/x/sys/windows/registry)**  
  License: BSD 3-Clause
  Copyright 2009 The Go Authors.
  Used for Windows registry access and manipulation.

All third-party copyright notices and licenses are preserved in this repository in compliance with their original terms.

## Contact

We'd love to hear from you! Here's how to get in touch:

### Issues (Bug Reports & Feature Requests)

Use [GitHub Issues](https://github.com/synapseq-foundation/synapseq/issues) for:

- Bug reports and technical problems
- Feature requests and enhancement suggestions
- Documentation improvements

### Discussions (Questions & Community)

Use [GitHub Discussions](https://github.com/synapseq-foundation/synapseq/discussions) for:

- General questions and support (e.g., "How do I use `@presetlist`?")
- Help with your sequences (e.g., "My sequence isn't working, can you help?")
- Sharing your own sequences and presets with the community
- Discussing ideas and best practices
- Showcasing creative use cases

### Quick Guidelines

- **Found a bug?** → Open an Issue
- **Want a new feature?** → Open an Issue
- **Need help or have questions?** → Start a Discussion
- **Want to share your sequences?** → Post in Discussions
- **General feedback or ideas?** → Start a Discussion

## Credits

Check out the [CREDITS](CREDITS.md) to see a list of all contributors and special thanks!
