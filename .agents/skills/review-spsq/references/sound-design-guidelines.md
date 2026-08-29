# SPSQ Sound and Composition Review

## Contents

- [Review stance](#review-stance)
- [Layer density](#layer-density)
- [Tonal and rhythmic interaction](#tonal-and-rhythmic-interaction)
- [Noise, effects, and stereo motion](#noise-effects-and-stereo-motion)
- [Temporal narrative](#temporal-narrative)
- [Responsible positioning](#responsible-positioning)
- [Limits of static analysis](#limits-of-static-analysis)
- [Optional rendered-audio review](#optional-rendered-audio-review)

## Review stance

Describe consequences before recommending changes. Distinguish:

- **fact**: directly declared or measured;
- **engine behavior**: established by parser, builder, or renderer;
- **inference**: likely intent suggested by names or comments;
- **preference**: one plausible artistic direction.

Do not define one ideal sound. A sparse meditation piece and a deliberately dense psychedelic piece require different judgments. Preserve the declared identity and ask what the complexity contributes.

## Layer density

For each effective preset, count:

- active tracks;
- tone-based tracks;
- simultaneous binaural, monaural, and isochronic methods;
- noise, ambiance, and music beds;
- animated tracks and effect types.

Several moderate amplitude values can still produce a dense mix. Do not sum amplitudes and call the result loudness, headroom, dB, or clipping.

Flag as a technical warning when layers are likely to compete for the same perceptual role, especially when:

- several tonal carriers occupy nearby ranges;
- several pulse methods compete simultaneously;
- broad noise or external audio may mask quiet tones;
- every layer moves and no stable reference remains;
- many state changes occur in a short interval.

Classify deliberate density as an artistic observation when comments and structure support it. Offer simplification only as an optional experiment.

## Tonal and rhythmic interaction

Review:

- nearby carriers that may produce complex beating or masking;
- multiple binaural beat rates whose stereo differences compete;
- simultaneous binaural, monaural, and isochronic rhythms;
- carrier movement combined with doppler movement;
- high contrast between consecutive beat rates;
- square and sawtooth waveforms, which contain stronger harmonics than sine;
- custom waveforms with steep segments, sharp corners, asymmetric phase origins, or abrupt shape changes, which may emphasize harmonics, aliasing, or awkward morphs;
- long exposure to prominent pulsing or bright carriers.

Do not claim a carrier is universally unsafe or uncomfortable from its number alone. Use cautious language such as “potentially prominent,” “may feel bright,” or “worth auditioning at low level.”

Binaural content is meaningfully evaluated with stereo headphones. Monaural and isochronic pulses are physically present in the signal and may sound more explicit.

## Noise, effects, and stereo motion

Noise character:

- white is brightest;
- pink is more perceptually balanced;
- brown emphasizes lower frequencies;
- noise `smooth` changes short-term roughness, not color.

Review whether noise stabilizes, masks, or overwhelms tonal detail. A high noise smoothness value does not guarantee softness if amplitude and other layers remain prominent.

Effects:

- `pan` moves stereo position;
- `modulation` varies amplitude;
- `doppler` adds pitch motion to tone tracks.
- `shift` gives external audio opposing frequency offsets from a mono-derived wet signal.

Multiple simultaneous pan movements can occupy the stereo field. Multiple modulation effects can make the mix breathe or pulse competitively. Doppler plus carrier/beat interpolation can create compound pitch motion. High shift intensity can reduce the source's original stereo image and color complex material; multiple shifted external tracks can make the spectrum and stereo field less stable.

Do not report simultaneous effects as wrong. Explain the likely complexity and whether a stable anchor remains.

## Temporal narrative

Treat the timeline as a sequence of intentions:

1. entrance;
2. establishment;
3. development;
4. peak or central state when relevant;
5. release;
6. landing.

Not every sequence needs all six. Judge them against comments, names, duration, and stated purpose.

Review:

- whether the opening is immediate or faded;
- whether phase lengths are proportionate to their role;
- whether the middle evolves or holds;
- whether contrast is perceptible but controlled;
- whether intensity rises without release;
- whether large parameter changes have enough time;
- whether steps reinforce or distract from the arc;
- whether the last active state has time to settle;
- whether final silence creates a fade or merely follows an already silent span.

A long static period may be monotonous, meditative, stabilizing, or deliberately suspended. State the duration and absence of parameter motion first; label the artistic meaning as inference.

## Responsible positioning

Inspect `#` and `##` comments, preset names, and accompanying user descriptions for claims of:

- cure, treatment, prevention, or diagnosis;
- guaranteed sleep or focus;
- guaranteed brain synchronization;
- guaranteed altered consciousness;
- reproduction of drug effects;
- specific neurological, psychological, or medical outcomes.

Recommend factual creative language: “designed for,” “intended as,” “explores,” or “may feel.” SynapSeq is a creative and experimental audio tool, not a medical device.

When relevant, keep advice brief:

- listen at a comfortable level;
- use stereo headphones for binaural material;
- do not listen during driving or attention-critical activity;
- stop if discomfort occurs.

Do not turn an ordinary report into a legal disclaimer.

## Limits of static analysis

Static SPSQ review cannot establish:

- perceived quality for every listener;
- real acoustic level or dB SPL;
- device, amplifier, or headphone behavior;
- clipping in rendered samples;
- actual content of uninspected ambiance/music;
- seamlessness of external loops;
- medical, cognitive, or brain-state outcomes.

Recommend low-volume critical listening when sound quality matters. Label every unmeasured conclusion accordingly.

## Optional rendered-audio review

Perform only on explicit request and after `-test` succeeds.

Render safely:

1. create a unique temporary directory;
2. choose an explicit WAV path inside it;
3. run `synapseq INPUT OUTPUT.wav` or the available repository fallback;
4. never rely on the default output name beside the source.

If available, use:

- `ffprobe` to measure duration, channel count, codec, and sample rate;
- `ffmpeg` `astats` for sample peaks and other signal statistics;
- `ffmpeg` `ebur128` for integrated and momentary loudness indicators.

Report exact commands and preserve tool units. A peak at or extremely close to full scale warrants investigation; it is not by itself proof that every listener will hear distortion. Loudness measurements are properties of the rendered file, not dB SPL at the ear.

Do not play audio automatically. Clean up only the exact temporary directory created for the review.
