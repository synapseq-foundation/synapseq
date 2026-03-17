<h1 align="center">SynapSeq</h1>

<p align="center">
  <p align="center">
  <a href="https://github.com/synapseq-foundation/synapseq/releases/latest"><img src="https://img.shields.io/github/v/release/synapseq-foundation/synapseq?color=blue&logo=github" alt="Release"></a>
  <a href="COPYING.txt"><img src="https://img.shields.io/badge/license-GPL%20v2-blue.svg?logo=open-source-initiative&logoColor=white" alt="License"></a>
  <a href="https://github.com/synapseq-foundation/synapseq/commits"><img src="https://img.shields.io/github/commit-activity/m/synapseq-foundation/synapseq?color=ff69b4&logo=git" alt="Commit Activity"></a>
</p>
</p>

<p align="center"><strong>Synapse-Sequenced Brainwave Generator</strong></p>

**SynapSeq** is a text-driven audio sequencer for building clear, repeatable brainwave and ambient sessions using a simple domain-specific language, written as SynapSeq sequences (.spsq).

For a local command reference and language overview in Markdown, see [USAGE](USAGE.md). You can also open the terminal manual with `synapseq -manual`.

Visit [synapseq.org](https://synapseq.com) for more information.

## Quick Start

The recommended way to install SynapSeq is through the platform package manager first. If that is not available for your environment, install the binary from the GitHub releases page.

### Linux and macOS

Install with Homebrew:

```bash
brew tap synapseq-foundation/synapseq
brew install synapseq
```

Or install from GitHub releases:

1. Download the archive for your platform from the latest GitHub release.
2. Extract the archive.
3. Move the `synapseq` binary to `/usr/local/bin`.
4. Make sure it is executable.

Example:

```bash
chmod +x synapseq
sudo mv synapseq /usr/local/bin/
synapseq -help
```

### Windows

Install with winget:

```powershell
winget update
winget install synapseq
```

Or install from GitHub releases:

1. Download the ZIP archive for Windows from the latest GitHub release.
2. Extract the archive.
3. Move `synapseq.exe` to a folder that is already in the system `%PATH%`, or add a dedicated folder for SynapSeq to `%PATH%`.
4. Open a new PowerShell or Command Prompt window and verify the installation.
5. Run `synapseq -install-file-association` to associate `.spsq` files with SynapSeq and enable additional Explorer context menu actions.

Example:

```powershell
synapseq -help
synapseq -install-file-association
```

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

- **[fatih/color](https://github.com/fatih/color)**  
  License: MIT  
  Used for colorized terminal output.

- **[beep](https://github.com/gopxl/beep)**  
  License: MIT  
  Used for audio encoding/decoding.

- **[golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys)**  
  License: BSD 3-Clause  
  Used for platform-specific system integration.

- **[go-colorable](https://github.com/mattn/go-colorable)**  
  License: MIT  
  Used indirectly for cross-platform ANSI color support.

- **[go-isatty](https://github.com/mattn/go-isatty)**  
  License: MIT  
  Used indirectly for terminal capability detection.

- **[pkg/errors](https://github.com/pkg/errors)**  
  License: BSD 2-Clause  
  Used indirectly for error wrapping and stack trace utilities.

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
