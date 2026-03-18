# SynapSeq Color Palette

This folder documents the shared SynapSeq color language used by the project's visual interfaces.

Even though SynapSeq is primarily a CLI project, the repository already includes browser-facing surfaces such as the HTML preview and the WASM example. This directory exists so those interfaces can reuse the same palette and maintain a consistent visual tone.

## Purpose

The material in this folder exists to:

- define SynapSeq's shared color palette explicitly
- align light and dark themes under the same warm visual identity
- reduce drift between preview, WASM, and future HTML interfaces
- document which tokens are shared with the Go package and which are CSS-only
- keep color usage consistent with the intended product tone

## What The Palette Should Convey

The SynapSeq palette should support both the contemplative and technical sides of the product.

- warmth
- calm focus
- technical clarity
- long-session readability
- atmosphere without decorative excess

In practice, that means warm neutrals, restrained contrast, terracotta emphasis, earthy semantic colors, and track-specific accents that remain readable without overwhelming the interface.

## Files In This Folder

- [pl-basic.md](/Users/ruanf/Dev/my/synapseq/pallete/pl-basic.md): main color palette reference
- [example.html](/Users/ruanf/Dev/my/synapseq/pallete/example.html): standalone palette/example page for local visual checks

## How To Use It

When creating or adjusting visual interfaces in the project:

- consult [pl-basic.md](/Users/ruanf/Dev/my/synapseq/pallete/pl-basic.md) first
- reuse existing tokens before introducing new colors
- keep dark mode as a continuation of the light theme, not a separate identity
- treat the Go package `internal/palette` as the shared core token set for terminal-facing code
- treat the CSS palette as the broader reference for web interfaces

## Scope

This folder is not a full design system and not a marketing brand guide.
It exists specifically to document and preserve the SynapSeq color palette used by interfaces in this repository.
