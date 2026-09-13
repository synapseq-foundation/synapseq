# SynapSeq Informal Grammar

This is a lightweight, line-oriented grammar intended to summarize the parser shape. It is not a full formal specification of every semantic validation rule.

```text
file                 = { line } ;

line                 = blank-line
                     | comment-line
                     | option-line
                     | preset-line
                     | track-line
                     | track-override-line
                     | timeline-line ;

blank-line           = whitespace-only ;

comment-line         = [indent] "#" text
                     | [indent] "##" text ;

option-line          = "@samplerate" integer
                     | "@volume" integer
                     | "@ambiance" name [path-or-url]
                     | "@music" name [path-or-url]
                     | "@waveform" name waveform-point waveform-point { waveform-point }
                     | "@transition" name transition-point transition-point { transition-point }
                     | "@extends" path-or-url ;

preset-line          = name
                     | name "from" name
                     | name "as" "template" ;

track-line           = indent2 tone-track
                     | indent2 noise-track
                     | indent2 ambiance-track
                     | indent2 music-track ;

tone-track           = [waveform-prefix] "tone" float tone-tail ;
tone-tail            = "amplitude" amplitude-value
                     | beat-kind float "amplitude" amplitude-value
                     | "effect" tone-effect float "intensity" float "amplitude" amplitude-value
                     | beat-kind float "effect" tone-effect float "intensity" float "amplitude" amplitude-value ;

noise-track          = "noise" noise-kind noise-tail ;
noise-tail           = "amplitude" amplitude-value
                     | "smooth" float "amplitude" amplitude-value
                     | "effect" noise-effect float "intensity" float "amplitude" amplitude-value
                     | "smooth" float "effect" noise-effect float "intensity" float "amplitude" amplitude-value ;

ambiance-track       = [waveform-prefix] "ambiance" name ambiance-tail ;
ambiance-tail        = "amplitude" amplitude-value
                     | "effect" ambiance-effect float "intensity" float "amplitude" amplitude-value ;

music-track          = [waveform-prefix] "music" name music-tail ;
music-tail           = "amplitude" amplitude-value
                     | "effect" music-effect float "intensity" float "amplitude" amplitude-value ;

waveform-prefix      = "waveform" waveform ;
waveform             = name ;
waveform-point       = float ;  (* 0 through 100; 2 through 16384 points *)
transition-point     = float ;  (* non-decreasing 0 through 100; 2 through 256 points; first 0, last 100 *)
amplitude-value      = float | "left" float "right" float ;  (* each 0 through 100 *)
beat-kind            = "binaural" | "monaural" | "isochronic" ;
noise-kind           = "white" | "pink" | "brown" ;
tone-effect          = "pan" | "modulation" | "doppler" | "shift" ;
noise-effect         = "pan" | "modulation" | "shift" ;
ambiance-effect      = "pan" | "modulation" | "doppler" | "shift" ;
music-effect         = "pan" | "modulation" | "doppler" | "shift" ;

track-override-line  = indent2 "track" track-index override-kind override-value ;
track-index          = integer ;
override-kind        = "tone"
                     | "binaural"
                     | "monaural"
                     | "isochronic"
                     | "waveform"
                     | "pan"
                     | "modulation"
                     | "doppler"
                     | "shift"
                     | "smooth"
                     | "amplitude"
                     | "left"
                     | "right"
                      | "intensity" ;
override-value       = signed-float | waveform ;

timeline-line        = time name [transition [steps]] ;
time                 = HH ":" MM ":" SS ;
transition           = builtin-transition | name ;
builtin-transition   = "steady" | "ease-in" | "ease-out" | "smooth" ;
steps                = integer ;

indent2              = exactly two leading spaces ;
name                 = validated identifier ;
integer              = strict base-10 integer ;
float                = strict decimal number ;
signed-float         = float with optional leading "+" or "-" ;
path-or-url          = local path without extension | remote URL ;
```

Use this grammar as a compact map of accepted line shapes. For semantic rules, timeline behavior, inheritance restrictions, path normalization, and complete validation guidance, see the [SPSQ documentation](https://synapseq.org/docs/spsq).
