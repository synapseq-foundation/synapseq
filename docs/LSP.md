# SPSQ Language Server

SynapSeq includes a Language Server Protocol (LSP) server for `.spsq` sequences and
`.spsc` preset collections. Start it through standard input and output:

```sh
synapseq -lsp
```

The server communicates only with LSP JSON-RPC messages on standard output. Do
not use it to render audio or pass a sequence filename in the same invocation.

## Capabilities

- Full-document synchronization for opened buffers.
- Parse and validation diagnostics while editing.
- Contextual completion for SPSQ options, tracks, rhythms, effects, transitions,
  waveforms, presets, and declared music or ambiance resources.

For documents backed by a local file URI, diagnostics also check local
`@extends`, `@ambiance`, and `@music` references. Untitled buffers are checked
without filesystem access. The LSP never downloads remote resources.

## Client configuration

Configure an LSP-capable editor to start `synapseq -lsp` for the `.spsq` and
`.spsc` extensions. The client should use normal LSP `stdio` framing and send
full document changes. No editor-specific extension is required by the server.
