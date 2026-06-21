# API Demos

This directory contains small programs demonstrating how to use SynapSeq as a
Go library.

## Binaural CLI

`binaural-cli` generates a WAV file containing a binaural beat. It demonstrates
how to:

- create a sequence with `spsq.New`;
- add a tone with `Tone`, `Binaural`, and `Amplitude`;
- build a timeline with optional fade-in and fade-out;
- load the sequence through the public `core` API;
- render the result as a WAV file.

### Run

From the repository root:

```bash
go run ./api-demo/binaural-cli
```

This generates `output.wav` with the default settings: a 440 Hz carrier,
10 Hz binaural beat, 20% amplitude, and 60-second duration.

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-c` | Carrier frequency in Hz | `440` |
| `-b` | Binaural beat frequency in Hz | `10` |
| `-a` | Amplitude percentage | `20` |
| `-i` | Fade-in duration in seconds | `0` |
| `-o` | Fade-out duration in seconds | `0` |
| `-d` | Total duration in seconds | `60` |
| `-f` | Output WAV file | `output.wav` |

### Examples

Generate a five-minute alpha-range session:

```bash
go run ./api-demo/binaural-cli \
  -c 220 \
  -b 10 \
  -a 20 \
  -d 300 \
  -f alpha.wav
```

Add 15-second fade-in and fade-out transitions:

```bash
go run ./api-demo/binaural-cli \
  -c 220 \
  -b 10 \
  -a 20 \
  -i 15 \
  -o 15 \
  -d 300 \
  -f alpha-fade.wav
```

The fade durations should be non-negative, and their sum should not exceed the
total duration.

### Build

To compile a standalone executable:

```bash
go build -o bin/binaural-cli ./api-demo/binaural-cli
```

Then run it directly:

```bash
./bin/binaural-cli -c 220 -b 10 -d 300 -f alpha.wav
```
