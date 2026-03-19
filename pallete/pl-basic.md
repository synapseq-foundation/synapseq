# SynapSeq Color Palette

## Overview

This document defines the SynapSeq color palette used across product interfaces.

It is intentionally limited to color tokens and color-usage rules.

The Go package `/internal/palette` currently defines only the shared core RGB tokens used by terminal-facing code.
This document is broader than the Go package because it also documents CSS-only tokens used by web interfaces.

Current implementation references:

- `/internal/preview/index.html`
- `/wasm/example.html`

This file is not a full design system specification.
It does not define layout, typography, spacing, components, or motion.

---

## Palette Intent

The SynapSeq palette should communicate:

- warmth
- calm focus
- technical clarity
- low visual aggression
- long-session readability

The palette is based on warm neutrals, terracotta emphasis, restrained semantic colors, and atmospheric track-color accents.

Dark mode should preserve the same emotional tone.
It should feel like the same product in lower light, not like a separate visual identity.

---

## Theme Model

SynapSeq uses a canonical warm light palette and a matching warm dark palette.

Recommended theme selector:

```css
html[data-theme="dark"] {
  color-scheme: dark;
}
```

If another selector is used, the token mapping should remain identical.

---

## Base Tokens

### Light Base Tokens

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

### Dark Base Tokens

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

### Token Reference

| Token            | Light value                          | Dark value                        | Usage                        |
| ---------------- | ------------------------------------ | --------------------------------- | ---------------------------- |
| `--bg`           | `#f4efe5`                            | `#161311`                         | main page background         |
| `--panel`        | `rgba(255, 252, 246, 0.86)`          | `rgba(31, 26, 23, 0.82)`          | standard translucent surface |
| `--panel-strong` | `rgba(255, 250, 241, 0.98)`          | `rgba(38, 32, 28, 0.96)`          | stronger inner surface       |
| `--text`         | `#201a15`                            | `#f3eadf`                         | primary text                 |
| `--muted`        | `#6b6259`                            | `#b8aa9a`                         | secondary text               |
| `--line`         | `rgba(32, 26, 21, 0.12)`             | `rgba(243, 234, 223, 0.12)`       | default border/separator     |
| `--line-strong`  | `rgba(32, 26, 21, 0.18)`             | `rgba(243, 234, 223, 0.18)`       | emphasized divider/border    |
| `--shadow`       | `0 24px 80px rgba(51, 37, 23, 0.16)` | `0 24px 80px rgba(0, 0, 0, 0.32)` | elevated surface shadow      |

---

## Accent and Semantic Tokens

### Light Accent and Semantic Tokens

```css
:root {
  --accent: #b14d2a;
  --accent-strong: #7f2d18;
  --accent-soft: rgba(177, 77, 42, 0.14);
  --ok: #2f6b45;
  --warn: #a07126;
  --danger: #8b2e2e;
}
```

### Dark Accent and Semantic Tokens

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

### Token Reference

| Token             | Light value               | Dark value                 | Usage                         |
| ----------------- | ------------------------- | -------------------------- | ----------------------------- |
| `--accent`        | `#b14d2a`                 | `#c96a42`                  | primary interaction emphasis  |
| `--accent-strong` | `#7f2d18`                 | `#efb08d`                  | stronger accent text/emphasis |
| `--accent-soft`   | `rgba(177, 77, 42, 0.14)` | `rgba(201, 106, 66, 0.18)` | soft accent tint              |
| `--ok`            | `#2f6b45`                 | `#58a06a`                  | success state                 |
| `--warn`          | `#a07126`                 | `#c89a46`                  | warning state                 |
| `--danger`        | `#8b2e2e`                 | `#cc6d6d`                  | error state                   |

### Shared Go Token Mapping

The following light-theme tokens are currently mirrored by `/internal/palette/pallete.go`:

| Go token         | Value     | CSS equivalent    |
| ---------------- | --------- | ----------------- |
| `Terracotta`     | `#b14d2a` | `--accent`        |
| `TerracottaDark` | `#7f2d18` | `--accent-strong` |
| `Ochre`          | `#a07126` | `--warn`          |
| `Green`          | `#2f6b45` | `--ok`            |
| `MutedWarm`      | `#6b6259` | `--muted`         |
| `DangerRed`      | `#8b2e2e` | `--danger`        |

Notes:

- The Go package currently exposes only shared core color tokens.
- Base surface tokens, shadow tokens, dark-theme variants, and track-type tokens are currently CSS-only.

---

## Track-Type Tokens

These tokens identify track categories in visualizations.

They are not part of the neutral UI surface palette, but they must remain consistent wherever track semantics are shown.

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

## Background Color System

### Light Background Composition

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

### Dark Background Composition

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
- Avoid cold grays, saturated blues, or neon-heavy dominant colors.
- Dark mode should feel like the same room with the lights lowered.

---

## Color Usage Rules

### Neutrals

- Use `--bg`, `--panel`, and `--panel-strong` for the main surface system.
- Use `--text` for primary content.
- Use `--muted` for secondary text and low-emphasis explanations.
- Use `--line` for default separators and borders.
- Use `--line-strong` only when additional structural definition is necessary.

### Accent

- Use `--accent` for primary actions and key highlights.
- Use `--accent-strong` for stronger text emphasis inside accent-led contexts.
- Use `--accent-soft` for subtle focus rings, hover tints, or selection backgrounds.

### Semantic States

- Use `--ok` for success and confirmation.
- Use `--warn` for warnings and caution states.
- Use `--danger` for errors and destructive states.
- Prefer translucent semantic fills over hard solid alert blocks.

### Track Colors

- Preserve track-color semantics consistently across charts, legends, and segment views.
- Do not reuse track colors as generic button or status colors unless the context is explicitly track-related.
- Track colors may be more saturated than the neutral UI palette, but they should still harmonize with the warm environment.

---

## Accessibility and Contrast

- Preserve readable contrast between `--text` and the surface palette.
- Do not reduce muted text or border visibility until structure becomes ambiguous.
- In dark mode, prefer warm off-white text instead of pure white to reduce glare.
- Color should not be the only semantic signal when another cue is feasible.

---

## Implementation Notes

- Tokens are currently defined in local `:root` blocks in the HTML implementations.
- The palette is CSS-first and does not require a separate token build pipeline.
- Reuse existing tokens before adding new ones.
- If a new color is introduced, it must fit the warm SynapSeq palette and be defined for both light and dark themes when applicable.

---

## Non-Goals

- This document does not define layout rules.
- This document does not define component specs.
- This document does not define typography, spacing, or motion.
- This document should not be used to justify a cold SaaS palette, neon accents, or cyberpunk dark mode.

---

## Consistency Rules

- Keep the palette warm, restrained, and readable.
- Use terracotta as the primary interaction accent.
- Keep semantic colors earthy rather than synthetic.
- Preserve the same emotional signature in both light and dark themes.
- Treat track-type colors as data semantics, not generic UI decoration.
