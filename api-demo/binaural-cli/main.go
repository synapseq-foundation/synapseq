package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	synapseq "github.com/synapseq-foundation/synapseq/v4/core"
	"github.com/synapseq-foundation/synapseq/v4/spsq"
)

func main() {
	carrier := flag.Float64("c", 440.0, "carrier frequency")
	beat := flag.Float64("b", 10.0, "binaural frequency")
	amplitude := flag.Float64("a", 20.0, "volume")
	fadeIn := flag.Int("i", 0, "fade in duration")
	fadeOut := flag.Int("o", 0, "fade out duration")
	duration := flag.Int("d", 60, "duration in seconds")
	file := flag.String("f", "output.wav", "output file")
	flag.Parse()

	fmt.Fprintln(os.Stderr, "Usage: binaural-cli -c <carrier> -b <beat> -a <amplitude> -i <fadeIn> -o <fadeOut> -d <duration> -f <file>")

	ctx := synapseq.NewAppContext()

	seq := spsq.New()
	ts := seq.NewPreset("tone-set").Tone(*carrier).Binaural(*beat).Amplitude(*amplitude)

	var tim *spsq.Builder
	if *fadeIn > 0 {
		tim = seq.SilenceAt(0).PresetAt(time.Duration(*fadeIn)*time.Second, ts)
	} else {
		tim = seq.PresetAt(0, ts)
	}

	diff := time.Duration(*duration)*time.Second - time.Duration(*fadeOut)*time.Second

	if *fadeOut > 0 {
		tim = tim.PresetAt(diff, ts).SilenceAt(time.Duration(*duration) * time.Second)
	} else {
		tim = tim.PresetAt(diff, ts)
	}

	loaded, err := tim.Load(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("rendering %s...\n", *file)
	if err := loaded.WAV(*file); err != nil {
		fmt.Println(err)
	}
	fmt.Println("done")
}
