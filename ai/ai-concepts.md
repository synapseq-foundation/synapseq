# SynapSeq AI Concepts Guide

This document complements the technical specification defined in `ai-basic.md`.
While the technical guide explains the syntax rules, this document explains the
concepts behind brainwave entrainment so AI models can generate more meaningful
sessions instead of random parameter combinations.

This guide should be used together with the technical specification.

This document does not replace the output contract from `ai-basic.md`.

When a model answers a user request:

- follow the output-format and out-of-scope rules from `ai-basic.md`
- use this file only to improve session design choices
- do not output conceptual explanations unless the calling system explicitly asks for explanation instead of code generation

---

## BRAINWAVE ENTRAINMENT BASICS

Brainwave entrainment attempts to influence brain activity using rhythmic stimuli.
In audio entrainment systems, this is typically done using pulses or frequency
differences that the brain can synchronize with over time.

The most common brainwave ranges are:

Delta: 0.5 – 4 Hz
Deep sleep and unconscious states.

Theta: 4 – 8 Hz
Meditation, creativity, deep relaxation, dream states.

Alpha: 8 – 12 Hz
Calm awareness, relaxed focus, light meditation.

Beta: 12 – 30 Hz
Active thinking, concentration, problem solving.

Sessions usually move gradually between states rather than jumping abruptly.

---

## BEAT TYPES

SynapSeq supports three beat generation methods.

BINAURAL BEATS

How it works:

Two slightly different frequencies are sent to each ear. The brain perceives
the difference between them as a rhythmic beat.

Example:

tone 200 binaural 10

This means:

Left ear ≈ 200 Hz
Right ear ≈ 210 Hz
Perceived beat = 10 Hz

Best use cases:

• relaxation
• meditation
• sleep preparation
• subtle entrainment

Characteristics:

• soft and subtle
• requires headphones
• commonly used for alpha and theta sessions

MONAURAL BEATS

How it works:

Two nearby frequencies are combined before reaching the ears, creating an
amplitude-modulated signal.

Characteristics:

• stronger than binaural
• works with speakers
• more noticeable pulse

Best use cases:

• moderate entrainment
• relaxation sessions
• background listening

ISOCHRONIC TONES

How it works:

A single tone is rapidly turned on and off at a specific rhythm.

Example:

tone 200 isochronic 12

This produces a clear rhythmic pulse.

Characteristics:

• very strong entrainment signal
• does not require headphones
• very noticeable

Best use cases:

• focus
• studying
• alertness
• cognitive stimulation

Isochronic tones are generally better for **higher frequencies**
such as beta or high alpha sessions.

---

## NOISE TYPES

Noise layers add texture and can make sessions more comfortable to listen to.

WHITE NOISE

Equal energy across all frequencies.

Characteristics:

• bright
• hiss-like
• wide spectrum

Common uses:

• masking background noise
• focus sessions
• stimulation sessions

PINK NOISE

Energy decreases gradually toward higher frequencies.

Characteristics:

• softer than white noise
• balanced sound
• very natural

Common uses:

• relaxation
• meditation
• general background ambience

Pink noise is often the safest default choice.

BROWN NOISE

Strong emphasis on low frequencies.

Characteristics:

• deep
• rumbling
• very smooth

Common uses:

• sleep sessions
• deep relaxation
• anxiety reduction

Brown noise is often preferred for long calming sessions.

---

## NOISE SMOOTHNESS

The `smooth` parameter controls how slowly the noise texture evolves over time.

Range:

0 – 100

Lower values:

• rougher noise
• more variation

Higher values:

• smoother sound
• slower spectral movement
• more ocean-like texture

Typical ranges:

20 – 40 : natural noise
40 – 60 : smooth relaxation sound
60+ : very calm evolving texture

---

## TONE CARRIER FREQUENCIES

The carrier frequency is the audible tone that carries the beat.

Example:

tone 200 binaural 8

Here:

200 Hz is the carrier
8 Hz is the entrainment beat

Typical carrier ranges:

120 – 300 Hz

Lower carriers sound deeper and warmer.
Higher carriers sound brighter and more energetic.

---

## WAVEFORMS

SynapSeq oscillators support multiple waveform shapes.
Each waveform has a different harmonic structure and perceptual character.

SINE

Characteristics:

• pure tone
• smooth
• no harmonics

Best uses:

• meditation
• relaxation
• long listening sessions

Sine waves are the safest default waveform.

TRIANGLE

Characteristics:

• soft harmonics
• warmer than sine
• slightly richer tone

Best uses:

• balanced entrainment
• calm focus sessions

SAWTOOTH

Characteristics:

• very rich harmonic spectrum
• bright and energetic
• can feel intense

Best uses:

• stimulation sessions
• advanced users

SQUARE

Characteristics:

• strong odd harmonics
• very sharp and aggressive tone

Important rule:

Square waves can sound harsh at higher amplitudes.

For **tone tracks**, the amplitude should ideally remain:

below 8%

Recommended safe range:

3 – 8 amplitude

Use cases:

• alertness sessions
• experimental stimulation

Square waves should generally be used cautiously and at low amplitudes.

---

## EFFECTS

SynapSeq supports spatial and movement effects that add variation
to otherwise static tones or noise layers.

PAN

Description:

Moves the sound left and right across the stereo field.

How it works:

The value controls the speed of movement.
The intensity controls how wide the movement becomes.

Best uses:

• subtle spatial motion
• ambiance layers
• meditation textures

MODULATION

Description:

Applies rhythmic amplitude movement.

How it works:

The value controls the modulation rate.
The intensity controls how strong the volume oscillation becomes.

Best uses:

• adding gentle movement to noise
• slow breathing-like effects
• preventing static soundscapes

DOPPLER

Description:

Simulates movement of a sound source relative to the listener.

This effect combines:

• stereo movement
• slight pitch shift

Characteristics:

• creates a drifting or orbiting sound sensation
• adds strong spatial dynamics

Best uses:

• tone tracks
• immersive sound design
• advanced sessions

Doppler should be used moderately because strong settings
can become distracting.

---

## SESSION STRUCTURE

Well-designed sessions usually follow a progression.

Example meditation structure:

1. Preparation phase
   Gradually introduce sound.

2. Descent phase
   Slowly move toward lower frequencies.

3. Deep phase
   Maintain the target entrainment state.

4. Return phase
   Gradually move back toward lighter states.

5. Exit phase
   Fade to silence.

Abrupt jumps should be avoided.

---

## TRANSITIONS

SynapSeq transitions define how the system moves between two timeline states.

STEADY

The default behavior.

Values remain constant until the next timeline point.

Use cases:

• stable segments
• holding a specific state

EASE-IN

Transition begins slowly and accelerates toward the target.

Use cases:

• session starts
• gentle introduction of sound

EASE-OUT

Transition begins quickly and slows down near the target.

Use cases:

• exiting a state
• calming endings

SMOOTH

Balanced easing curve across the entire transition.

Use cases:

• gradual meditation ramps
• most relaxation transitions

Smooth is usually the safest default for gradual state changes.

---

## SESSION DESIGN GUIDELINES

Good AI-generated sessions should follow these principles:

1. Avoid abrupt frequency jumps.
2. Move gradually between brainwave ranges.
3. Keep carrier frequencies relatively stable.
4. Maintain consistent track structure between presets.
5. Use noise to soften the soundscape.
6. Use transitions to create natural progression.

---

## IMPORTANT NOTE FOR AI MODELS

This document explains concepts only.

For valid `.spsq` syntax generation, always follow the rules defined in:

ai-basic.md
