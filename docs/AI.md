# AI Tools

SynapSeq supports SPSQ generation through OpenAI-compatible APIs and provides
three skills for AI coding agents.

## Generate SPSQ From The CLI

Use `-ai` with an OpenAI-compatible API:

```bash
export SYNAPSEQ_AI_API_KEY="your-api-key"
synapseq -ai "Generate a 10 minute relaxation sequence"
```

Choose an output path, or use `-` to write only SPSQ content to standard
output:

```bash
synapseq -ai "Generate a 30 minute study sequence" study-30m.spsq
synapseq -ai "Generate a 20 minute meditation sequence" -
```

SynapSeq validates every model response before writing it. Invalid responses
are repaired automatically up to two times; if validation still fails, no
output file is created. Existing destination files are never overwritten.

### Configuration

The default model is `gpt-4.1-mini`, the default API host is OpenAI, the CLI
temperature is `1`, and requests time out after five minutes.

```bash
export SYNAPSEQ_AI_API_KEY="local-key"
export SYNAPSEQ_AI_MODEL="google/gemma-4-e4b"
export SYNAPSEQ_AI_BASE_URL="http://localhost:1234"
export SYNAPSEQ_AI_TEMPERATURE="0.2"
export SYNAPSEQ_AI_TIMEOUT="90s"
synapseq -ai "Generate a 15 minute focus sequence"
```

The same settings can be provided as CLI flags:

```text
-ai-model MODEL
-ai-base-url URL
-ai-temperature VALUE
-ai-timeout DURATION
```

On Windows, set the variables in PowerShell:

```powershell
$env:SYNAPSEQ_AI_API_KEY = "your-api-key"
$env:SYNAPSEQ_AI_MODEL = "gpt-4.1-mini"
synapseq -ai "Generate a 15 minute focus sequence"
```

Or in Command Prompt:

```bat
set SYNAPSEQ_AI_API_KEY=your-api-key
set SYNAPSEQ_AI_MODEL=gpt-4.1-mini
synapseq -ai "Generate a 15 minute focus sequence"
```

For English prompts containing terms such as `delta`, `theta`, `alpha`,
`beta`, `gamma`, `sleep`, `meditation`, `focus`, or `relaxation`, SynapSeq
also validates beat ranges, audible carriers, and required profile
progression. These are sound-design constraints, not promises about listener
outcomes.

For the equivalent Go API, see the [Go API Guide](GO_API.md#ai-generation-from-go).

## AI Agent Skills

SynapSeq publishes three complementary
[skills](https://skills.sh/synapseq-foundation/synapseq) for AI coding agents:

| Skill | Use it to | Result | Existing files |
|-------|-----------|--------|----------------|
| [`create-spsq`](https://skills.sh/synapseq-foundation/synapseq/create-spsq) | Create a sequence or derive a new version from a reference | A new validated `.spsq` file | Never modified or overwritten |
| [`explain-spsq`](https://skills.sh/synapseq-foundation/synapseq/explain-spsq) | Learn syntax or understand an existing sequence | A didactic, read-only explanation | Never modified |
| [`review-spsq`](https://skills.sh/synapseq-foundation/synapseq/review-spsq) | Audit syntax, semantics, sound structure, and timeline composition | A prioritized read-only report | Never modified |

Choose the skill by the main action: **create**, **explain**, or **review**.
Only `create-spsq` writes complete sequence files, and it always writes to a
new path.

### Install

Install the complete skill suite with the [`skills` CLI](https://github.com/vercel-labs/skills):

```bash
npx skills add synapseq-foundation/synapseq
```

The command lets you choose the target agent and installation scope.

### Codex

Mention the skill explicitly:

```text
$create-spsq Create a 20-minute relaxation sequence with a smooth binaural fade-in and fade-out.
$explain-spsq Explain the presets, tracks, and timeline in focus.spsq.
$review-spsq Audit focus.spsq and report technical, structural, and artistic findings.
```

Codex may also select a skill automatically when the request matches its
description.

### Claude Code

Invoke an installed skill as a slash command:

```text
/create-spsq Create a 20-minute relaxation sequence with a smooth binaural fade-in and fade-out.
/explain-spsq Explain the presets, tracks, and timeline in focus.spsq.
/review-spsq Audit focus.spsq and report technical, structural, and artistic findings.
```

Claude Code can also load a skill automatically. Explicit invocation is the
most predictable option. Install SynapSeq using the [Quick Start](../README.md#quick-start)
instructions when the agent needs to validate a local file with `synapseq -test`.
