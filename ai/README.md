# SynapSeq AI Specification Pack

This folder contains the LLM-oriented specification files for generating SynapSeq `.spsq` sequences.

The documents are intentionally split by responsibility so a model can distinguish:

- strict syntax and response behavior
- transition compatibility constraints
- session design and entrainment concepts

## Files

### ai-basic.md

Use this as the primary document.

It defines:

- the required output contract
- the out-of-scope fallback response
- the minimal safe SynapSeq syntax
- timeline basics
- transition basics
- silence semantics

If a model receives only one file from this folder, it should be this one.

### ai-compatibility.md

Use this together with `ai-basic.md` when the model must generate multi-stage sessions with smooth transitions.

It defines:

- structural compatibility between directly connected presets
- track-position consistency
- beat-type compatibility
- noise-type compatibility
- effect-type compatibility
- when `silence` must be used as a reset bridge

This file does not replace the output contract from `ai-basic.md`.

### ai-concepts.md

Use this together with `ai-basic.md` when the model should generate sessions that are not only valid, but also meaningful from a brainwave-entrainment perspective.

It defines:

- entrainment ranges such as delta, theta, alpha, and beta
- practical differences between binaural, monaural, and isochronic beats
- noise and waveform characteristics
- effect usage guidance
- session design patterns and progression ideas

This file is conceptual guidance, not a syntax contract.

## Recommended Usage Order

For most systems, load the files in this order:

1. `ai-basic.md`
2. `ai-compatibility.md`
3. `ai-concepts.md`

Reason:

- `ai-basic.md` defines what the model is allowed to output
- `ai-compatibility.md` prevents structurally invalid transitions
- `ai-concepts.md` improves the quality of the generated session design

## Suggested Loading Strategies

### Minimal safe generation

Load:

- `ai-basic.md`

Use this when the goal is:

- short sessions
- simple presets
- lowest possible risk of parser-invalid output

### Safe multi-stage generation

Load:

- `ai-basic.md`
- `ai-compatibility.md`

Use this when the goal is:

- multiple presets
- gradual transitions
- structurally safe timeline generation

### Quality-oriented generation

Load:

- `ai-basic.md`
- `ai-compatibility.md`
- `ai-concepts.md`

Use this when the goal is:

- full meditation or focus sessions
- gradual state progression
- more realistic entrainment design choices

## Operational Rule

When documents overlap, follow this priority:

1. `ai-basic.md`
2. `ai-compatibility.md`
3. `ai-concepts.md`

In practice:

- if `ai-basic.md` defines an output restriction, that restriction wins
- if `ai-compatibility.md` forbids a direct structural transition, do not override it with conceptual guidance from `ai-concepts.md`
- use `ai-concepts.md` only to improve choices inside the limits defined by the other two files

## Intended Outcome

A well-configured model using this folder should:

- return only `.spsq` code when the request is in scope
- reject out-of-scope requests using the fallback defined in `ai-basic.md`
- avoid parser-invalid syntax
- avoid incompatible direct transitions
- produce more coherent entrainment session designs
