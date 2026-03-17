# SynapSeq Design System

## Overview

This document defines the implementation-level UI specification used across SynapSeq interfaces, currently centered on:

- `/internal/preview/index.html`
- `/wasm/example.html`

The purpose of this file is to provide a technical reference for tokens, layout rules, components, interaction patterns, and implementation constraints so future SynapSeq interfaces remain visually consistent.

---

## Design Intent

The SynapSeq interface should communicate the following qualities:

- contemplative and immersive
- technical and precise
- calm rather than aggressive
- structured without feeling clinical
- readable for long-form analytical use

The visual language is based on a warm light theme with translucent surfaces, restrained contrast, generous rounding, and technical typography for measurements and code.

Dark mode must preserve those same qualities. It should feel nocturnal and focused, not cyberpunk, neon, or clinical.

---

## Theme Application

### Canonical Theme

SynapSeq standardizes a canonical warm light theme and a matching warm dark theme.

- Background: warm cream gradient with low-opacity teal and amber radial layers
- Cards and panels: off-white translucent surfaces with blur
- Borders: warm low-alpha brown/ink lines
- Text primary: deep warm ink
- Text secondary: muted warm gray-brown
- Accent: terracotta
- Semantic states: green, ochre, earthy red

### Dark Mode

Dark mode is part of the design system when a darker viewing environment is needed, but it must preserve the same emotional tone as the light theme.

- Background: deep warm-charcoal gradient with restrained teal and amber atmospheric layers
- Cards and panels: charcoal translucent surfaces with blur, not flat black slabs
- Borders: soft warm gray lines with slightly stronger definition than light mode
- Text primary: warm off-white rather than cold pure white
- Text secondary: muted warm gray
- Accent: terracotta family, slightly lifted for contrast
- Semantic states: green, ochre, and red variants tuned for dark surfaces

### Theme Switching Model

Recommended implementation:

```css
html[data-theme="dark"] {
  color-scheme: dark;
}
```

If another selector is preferred, the token mapping should remain identical.

---

## Color Palette

### Base Tokens

```css
:root {
  --bg: #f4efe5;
  --panel: rgba(255, 252, 246, 0.86);
  --panel-strong: rgba(255, 250, 241, 0.98);
  --text: #201a15;
  --muted: #6b6259;
  --line: rgba(32, 26, 21, 0.12);
  --line-strong: rgba(32, 26, 21, 0.18);
  --shadow: 0 24px 80px rgba(51, 37, 23, 0.16);
}
```

```css
html[data-theme="dark"] {
  --bg: #161311;
  --panel: rgba(31, 26, 23, 0.82);
  --panel-strong: rgba(38, 32, 28, 0.96);
  --text: #f3eadf;
  --muted: #b8aa9a;
  --line: rgba(243, 234, 223, 0.12);
  --line-strong: rgba(243, 234, 223, 0.18);
  --shadow: 0 24px 80px rgba(0, 0, 0, 0.32);
}
```

| Token            | Value                                | Usage                           |
| ---------------- | ------------------------------------ | ------------------------------- |
| `--bg`           | `#f4efe5`                            | main page background            |
| `--panel`        | `rgba(255, 252, 246, 0.86)`          | standard translucent surface    |
| `--panel-strong` | `rgba(255, 250, 241, 0.98)`          | stronger card/input surface     |
| `--text`         | `#201a15`                            | primary text                    |
| `--muted`        | `#6b6259`                            | secondary text                  |
| `--line`         | `rgba(32, 26, 21, 0.12)`             | default border/separator        |
| `--line-strong`  | `rgba(32, 26, 21, 0.18)`             | emphasized border or guide line |
| `--shadow`       | `0 24px 80px rgba(51, 37, 23, 0.16)` | default elevated shadow         |

#### Dark Base Token Mapping

| Token            | Value                             | Usage                           |
| ---------------- | --------------------------------- | ------------------------------- |
| `--bg`           | `#161311`                         | main page background            |
| `--panel`        | `rgba(31, 26, 23, 0.82)`          | standard translucent surface    |
| `--panel-strong` | `rgba(38, 32, 28, 0.96)`          | stronger card/input surface     |
| `--text`         | `#f3eadf`                         | primary text                    |
| `--muted`        | `#b8aa9a`                         | secondary text                  |
| `--line`         | `rgba(243, 234, 223, 0.12)`       | default border/separator        |
| `--line-strong`  | `rgba(243, 234, 223, 0.18)`       | emphasized border or guide line |
| `--shadow`       | `0 24px 80px rgba(0, 0, 0, 0.32)` | default elevated shadow         |

### Accent and Semantic Tokens

```css
:root {
  --accent: #b14d2a;
  --accent-strong: #7f2d18;
  --accent-soft: rgba(177, 77, 42, 0.14);
  --ok: #2f6b45;
  --warn: #885b17;
  --danger: #8b2e2e;
}
```

```css
html[data-theme="dark"] {
  --accent: #c96a42;
  --accent-strong: #efb08d;
  --accent-soft: rgba(201, 106, 66, 0.18);
  --ok: #58a06a;
  --warn: #c89a46;
  --danger: #cc6d6d;
}
```

| Token             | Value                     | Usage                        |
| ----------------- | ------------------------- | ---------------------------- |
| `--accent`        | `#b14d2a`                 | primary action background    |
| `--accent-strong` | `#7f2d18`                 | accent text, active emphasis |
| `--accent-soft`   | `rgba(177, 77, 42, 0.14)` | hover/focus/soft highlight   |
| `--ok`            | `#2f6b45`                 | success state                |
| `--warn`          | `#885b17`                 | warning state                |
| `--danger`        | `#8b2e2e`                 | error state                  |

#### Dark Accent and Semantic Token Mapping

| Token             | Value                      | Usage                        |
| ----------------- | -------------------------- | ---------------------------- |
| `--accent`        | `#c96a42`                  | primary action background    |
| `--accent-strong` | `#efb08d`                  | accent text, active emphasis |
| `--accent-soft`   | `rgba(201, 106, 66, 0.18)` | hover/focus/soft highlight   |
| `--ok`            | `#58a06a`                  | success state                |
| `--warn`          | `#c89a46`                  | warning state                |
| `--danger`        | `#cc6d6d`                  | error state                  |

### Track-Type Tokens

```css
:root {
  --pure: #4f46e5;
  --binaural: #0f766e;
  --monaural: #1d4ed8;
  --isochronic: #b45309;
  --noise: #9a3412;
  --ambiance: #047857;
  --silence: #6b7280;
  --off: #94a3b8;
}
```

| Token          | Value     | Usage                         |
| -------------- | --------- | ----------------------------- |
| `--pure`       | `#4f46e5` | pure tone visualization       |
| `--binaural`   | `#0f766e` | binaural beat visualization   |
| `--monaural`   | `#1d4ed8` | monaural beat visualization   |
| `--isochronic` | `#b45309` | isochronic beat visualization |
| `--noise`      | `#9a3412` | noise visualization           |
| `--ambiance`   | `#047857` | ambiance visualization        |
| `--silence`    | `#6b7280` | silence/inactive segment      |
| `--off`        | `#94a3b8` | disabled/off state            |

---

## Background System

### Canonical Page Background

```css
background:
  radial-gradient(
    circle at top left,
    rgba(15, 118, 110, 0.18),
    transparent 24%
  ),
  radial-gradient(circle at top right, rgba(180, 83, 9, 0.16), transparent 20%),
  linear-gradient(180deg, #fbf6eb 0%, var(--bg) 100%);
```

### Canonical Dark Background

```css
background:
  radial-gradient(
    circle at top left,
    rgba(15, 118, 110, 0.14),
    transparent 26%
  ),
  radial-gradient(circle at top right, rgba(180, 83, 9, 0.14), transparent 22%),
  linear-gradient(180deg, #1c1714 0%, var(--bg) 100%);
```

### Rules

- Never use a flat solid page background for primary SynapSeq surfaces.
- Radial gradients should remain low-opacity and atmospheric.
- The background must support long reading sessions and dense interface content.
- Avoid cold grays, saturated blues, or high-energy neon accents as the dominant base.
- Dark mode should feel like the same room with the lights lowered, not like a different product.

---

## Typography

### Font Families

- Primary UI: `"Avenir Next", "Segoe UI", sans-serif`
- Monospace: `Menlo, Monaco, "Cascadia Mono", "SFMono-Regular", monospace`

### Font Scale

| Element               | Size             |
| --------------------- | ---------------- |
| Eyebrow / micro label | `12px`           |
| Secondary UI text     | `13px` to `14px` |
| Body text             | `16px` to `17px` |
| Section title         | `24px` to `34px` |
| Hero title            | `36px` to `68px` |
| Highlight value       | `22px` to `34px` |

### Font Weights

- Regular: `400`
- Medium: `500`
- Semibold: `600`
- Bold: `700`

### Typographic Rules

- Eyebrows and small labels use uppercase with high letter-spacing.
- Major headings use negative tracking for a compact editorial feel.
- Body copy should stay soft and readable, using `--muted` where appropriate.
- All technical values, code, timestamps, and measurement outputs use the monospace stack.

---

## Spacing Scale

### Core Spacing Values

| Tokenized usage      | Size             |
| -------------------- | ---------------- |
| Micro gap            | `8px`            |
| Tight control gap    | `10px`           |
| Standard control gap | `12px`           |
| Small block gap      | `14px`           |
| Compact section gap  | `16px`           |
| Medium inner padding | `18px`           |
| Standard section gap | `20px` to `24px` |
| Large panel padding  | `28px`           |

### Container Widths

```css
/* preview */
width: min(1360px, calc(100vw - 40px));

/* wasm/demo */
width: min(1100px, calc(100vw - 40px));
```

---

## Border Radius

| Use                        | Radius           |
| -------------------------- | ---------------- |
| Input / compact block      | `18px`           |
| Card / stat / medium block | `20px` to `22px` |
| Large panel / shell        | `24px` to `28px` |
| Pill / chip / tab          | `999px`          |

Rule: SynapSeq surfaces should feel soft and tactile; avoid sharp-cornered primary components.

---

## Shadows and Depth

### Standard Shadow

```css
box-shadow: 0 24px 80px rgba(51, 37, 23, 0.16);
```

### Dark Shadow

```css
box-shadow: 0 24px 80px rgba(0, 0, 0, 0.32);
```

### Depth Rules

- Shadows must be large and diffuse, not hard-edged.
- Shadow color should remain warm-brown rather than pure black.
- Depth is created through blur, translucency, and layering rather than stark contrast.

---

## Surface Specification

### Base Panel

```css
background: var(--panel);
backdrop-filter: blur(14px);
border: 1px solid var(--line);
border-radius: 28px;
box-shadow: var(--shadow);
```

### Dark Base Panel

```css
background: var(--panel);
backdrop-filter: blur(14px);
border: 1px solid var(--line);
border-radius: 28px;
box-shadow: var(--shadow);
```

### Strong Inner Surface

```css
background: var(--panel-strong);
border: 1px solid rgba(32, 26, 21, 0.08);
border-radius: 20px;
```

### Dark Strong Inner Surface

```css
background: var(--panel-strong);
border: 1px solid rgba(243, 234, 223, 0.08);
border-radius: 20px;
```

### Usage

- Use `--panel` for structural shells and major containers.
- Use `--panel-strong` for inset cards, stats, editors, and focused content blocks.
- Use `--line` for default borders and `--line-strong` only when the guide needs stronger definition.
- In dark mode, preserve translucency and warmth; avoid pure black panels and stark white borders.

---

## Component Specifications

### Navbar

**Structure:**

- left-aligned brand
- muted subtitle adjacent to brand
- single rounded external link on the right

**Behavior:**

- compact height
- translucent panel treatment
- hover state uses soft accent background and slight upward movement
- dark mode should preserve subtlety; the navbar should separate through contrast, not brightness

### Hero

**Structure:**

- eyebrow label
- large title
- supporting paragraph with restrained line length
- adjacent or secondary stats block

**Rules:**

- hero titles use high visual prominence but restrained color
- supporting text should remain in `--muted`
- line length should remain readable for documentation and analytical contexts

### Stats Cards

**Base:**

- `--panel-strong` background
- subtle border
- uppercase micro-label
- oversized numeric or key value

**Usage:**

- sequence duration
- segment count
- track count
- progress or runtime summary

### Tabs and Pills

**Base:**

- full pill radius
- quiet neutral background
- compact horizontal padding

**Active State:**

- accent-led background or emphasis
- readable contrast

**Hover State:**

- `translateY(-1px)` max
- optional `--accent-soft` tint

**Dark Mode Notes:**

- keep pill backgrounds close to the panel family, not pitch-black
- active states should use accent or raised contrast without becoming luminous neon

### Graph UI

**Base Characteristics:**

- lightly glassy graph shell
- mono labels for rulers and values
- rounded SVG stroke joins and caps
- legends rendered as pills
- overlays and graph series must share the same time axis

**Rules:**

- metric tabs may change legend semantics
- legends can be track-based or category-based depending on the graph metric
- technical annotation must be brief and visually secondary to the graph lines

### Track and Segment Cards

**Base:**

- strong light surface
- soft border
- clear title hierarchy
- metadata chips
- muted summary
- compact data grid for technical values

**Usage:**

- track summaries
- node or segment metadata
- effect, waveform, noise, and channel summaries

### Input and Text Editing

**Base:**

```css
background: var(--panel-strong);
border: 1px solid var(--line);
border-radius: 18px;
font:
  14px/1.55 Menlo,
  Monaco,
  "Cascadia Mono",
  "SFMono-Regular",
  monospace;
```

**Focus:**

```css
outline: 2px solid rgba(177, 77, 42, 0.2);
border-color: rgba(177, 77, 42, 0.28);
```

**Dark Focus:**

```css
outline: 2px solid rgba(201, 106, 66, 0.24);
border-color: rgba(201, 106, 66, 0.34);
```

### Buttons

**Primary Button:**

- background: `--accent`
- text: light warm white
- shape: pill
- font weight: `600`
- hover transform: `translateY(-1px)`

**Secondary Button:**

- quiet neutral background
- subtle border
- warm text color

**Disabled State:**

- reduced opacity
- no transform
- non-interactive cursor

**Dark Mode Notes:**

- primary buttons may become slightly brighter than light mode to preserve emphasis
- secondary buttons should remain quiet and panel-adjacent, not glossy or metallic

### Status and Feedback

**Success:** use `--ok`

**Warning:** use `--warn`

**Error:** use `--danger`

**Style Rules:**

- semantic blocks should use translucent fills rather than solid alerts
- border and text should align with the semantic color family
- message hierarchy should remain readable inside dense operational UIs
- in dark mode, semantic fills should remain muted enough to avoid visual noise accumulation

---

## Motion and Transitions

### Standard Transition

```css
transition:
  transform 140ms ease,
  background 140ms ease,
  border-color 140ms ease,
  opacity 140ms ease;
```

### Motion Rules

- Interaction motion must be short and functional.
- Avoid decorative long-duration animation.
- Use movement primarily to confirm interaction, not to entertain.
- Standard hover lift should stay between `-1px` and `-2px`.

---

## Responsive Rules

### Layout Behavior

- Multi-column layouts collapse to a single column below `900px` or `720px`, depending on density.
- Horizontal control rows may stack vertically on smaller viewports.
- Shell margins tighten for small screens, especially below `640px`.

### Responsive Principle

Mobile layouts should preserve clarity and atmosphere without compressing typography or overloading horizontal space.

---

## Accessibility Guidelines

### Readability

- Body text should remain at `16px` or larger in reading-heavy contexts.
- Secondary UI text should generally not drop below `13px`.
- Prose blocks should stay near a readable line length.

### Interaction

- Interactive elements must expose visible hover and focus states.
- Focus treatments should remain consistent with the terracotta accent family.
- Color should not be the only indicator of active or semantic state when additional cues are feasible.

### Contrast

- Preserve readable contrast between `--text` and the warm panel surfaces.
- Do not reduce borders or muted text to the point that structure disappears on lower-quality displays.
- In dark mode, prefer warm off-white text over pure white to reduce glare during long sessions.

---

## Implementation Notes

### Current Sources of Truth

- `/internal/preview/index.html`
- `/wasm/example.html`

### Implementation Model

- Tokens are currently defined in local `:root` blocks inside the HTML files.
- The system is CSS-first and does not depend on an external design-token build pipeline.
- Reuse the existing tokens before introducing new values.

### Non-Goals

- This specification does not define a separate corporate marketing style distinct from the product UI.
- This specification should not be used to justify introducing generic blue SaaS styling that conflicts with the established SynapSeq tone.
- This specification should not be used to justify a cold, neon, cyberpunk dark theme.

---

## Consistency Rules

- Keep the background warm, layered, and atmospheric.
- Keep panels translucent and softly elevated.
- Prefer pills, cards, and rounded shells over rigid boxed components.
- Use monospace only where technical precision benefits readability.
- Preserve track-color semantics consistently across visualizations.
- Use terracotta for primary interaction emphasis.
- Avoid cold neutral palettes, sharp corners, and heavy black shadows.
- Ensure dark mode preserves the same emotional signature as light mode: calm, immersive, technical, and warm.

---

## Version History

- `v1.1.0` - Reframed the design system as a technical implementation specification and aligned it with the unified preview/WASM visual language.
- `v1.2.0` - Added canonical dark mode guidance, token mappings, and component behavior rules consistent with the light theme emotional tone.
