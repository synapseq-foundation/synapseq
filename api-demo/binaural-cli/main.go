package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/external"
	"github.com/synapseq-foundation/synapseq/v4/spsq"
)

func main() {
	carrier := flag.Float64("c", 440.0, "carrier frequency")
	beat := flag.Float64("b", 10.0, "binaural frequency")
	amplitude := flag.Float64("a", 20.0, "volume")
	fadeIn := flag.Int("i", 0, "fade in duration")
	fadeOut := flag.Int("o", 0, "fade out duration")
	duration := flag.Int("d", 60, "duration in seconds")
	help := flag.Bool("h", false, "show help")
	file := flag.String("f", "", "output file. Supported formats: .wav, .mp3")
	flag.Parse()

	if *help {
		fmt.Fprintln(os.Stderr, "Usage: binaural-cli [options]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Create a synapseq context
	ctx := synapseq.NewAppContext()

	// Create a new spsq sequence
	seq := spsq.New()
	// Create a new preset
	ts := seq.NewPreset("tone-set").Tone(*carrier).Binaural(*beat).Amplitude(*amplitude)

	// Create a timeline builder
	var tim *spsq.Builder
	// Initialize timeline with fade in if greater than 0, otherwise start with the preset
	if *fadeIn > 0 {
		tim = seq.SilenceAt(0).PresetAt(time.Duration(*fadeIn)*time.Second, ts)
	} else {
		tim = seq.PresetAt(0, ts)
	}

	// Calculate the duration of the sequence
	// Check if fade out is greater than 0 and adjust duration accordingly
	diff := time.Duration(*duration)*time.Second - time.Duration(*fadeOut)*time.Second

	if *fadeOut > 0 {
		tim = tim.PresetAt(diff, ts).SilenceAt(time.Duration(*duration) * time.Second)
	} else {
		tim = tim.PresetAt(diff, ts)
	}

	// Load the sequence from the timeline builder
	loaded, err := tim.Load(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid sequence\n"+err.Error())
		os.Exit(1)
	}

	// Save the sequence to a file if specified, otherwise play it using ffplay
	if *file != "" {
		ext := filepath.Ext(*file)
		if ext == ".mp3" {
			ffmpeg, err := external.NewFFmpeg("")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Fprint(os.Stderr, "Converting to MP3...")
			if err := ffmpeg.Convert(loaded, *file); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		} else {
			if err := loaded.WAV(*file); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		fmt.Fprintln(os.Stderr, "done")
		return
	}

	ffplay, err := external.NewFFPlay("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "playing...")
	if err := ffplay.Play(loaded); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "done")
}
