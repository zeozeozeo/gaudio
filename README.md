# Gaudio

Pure Go audio library with a simple and easy to use API (inspired by [pydub](https://github.com/jiaaro/pydub))

```go
import "github.com/zeozeozeo/gaudio"
```

# Supported formats

Audio formats:

-   MP3 ([go-mp3](https://github.com/hajimehoshi/go-mp3))
-   FLAC ([flac](https://github.com/mewkiz/flac))
-   WAV ([go-dsp/wav](https://github.com/mjibson/go-dsp/wav) and [go-wav](github.com/youpy/go-wav))

Tracker formats:

-   MOD ([gomodplay](https://github.com/zeozeozeo/gomodplay) and [Gotracker](https://github.com/gotracker/playback))
-   XM ([Gotracker](https://github.com/gotracker/playback))
-   S3M ([Gotracker](https://github.com/gotracker/playback))
-   IT ([Gotracker](https://github.com/gotracker/playback))

# Examples

## Create a segment from path

```go
import (
    "github.com/zeozeozeo/gaudio"
)

// Create a segment from path (see supported formats)
segment, err := gaudio.LoadAudioFromPath("audio.mp3", gaudio.FormatMP3)
if err != nil {
    panic(err)
}
```

## Create a segment from a reader

```go
// Open file
file, err := os.Open("song.mp3")
if err != nil {
    panic(err)
}
defer file.Close()

// Create a segment (see supported formats)
segment, err := gaudio.LoadAudio(file, gaudio.FormatMP3)
if err != nil {
    panic(err)
}
```

## Overlay segment2 over segment1 at the 3rd second

```go
segment1.Overlay(segment2, 3)
```

## Clone a segment

```go
clone := segment.Clone()
```

## Concatenate segments

```go
// Concatenate segment at the start of another segment
segment1.ConcatenateStart(segment2)

// Concatenate segment at the end of another segment
segment1.ConcatenateEnd(segment2)

// Concatenate segment at the specified second of another segment (3.5th second)
segment1.ConcatenateAtSecond(segment2, 3.5)
```

## Detect and remove silence

```go
threshold = 0.01 // Silence threshold

// Get leading and trailing silence
leading, trailing := segment.DetectSilence(threshold)

// Get leading silence
leading := segment.DetectSilenceStart(threshold)

// Get trailing silence
trailing := segment.DetectSilenceEnd(threshold)

// Remove silence
segment.RemoveStartAndEndSilence(threshold)
segment.RemoveStartSilence(threshold)
segment.RemoveEndSilence(threshold)
```

## Calculate BPM (warning: can be slow)

```go
// minBpm: the minimum BPM you expect (120)
// maxBpm: the maximum BPM you expect (200)
// steps: the amount of steps per interval (1024)
// samplesPerBeat: the amount of samples used for a single beat (1024)
bpm := segment.CalcBPM(120, 200, 1024, 1024)
```

## Apply echo effect

```go
// delay: delay between echo layers (0.05)
// layers: amount of layer (2)
segment.EchoEffect(0.05, 2)
```

## Pitch to note note (needs testing)

```go
// Pitch formula: 2 ** (key / 12)
segment.ApplyNote(12)
```

## Export a segment to a file (only wav files are supported right now)

```go
file, err := os.Create("output.wav")
if err != nil {
    panic(err)
}
defer file.Close()

segment.Export(file, gaudio.FormatWAVE)
```

## Apply fade in / fade out effect

```go
// start float64: the amount of fade-in seconds
// end float64: the amount of fade-out seconds
segment.Fade(0.5, 1) // Fade in for 0.5 seconds, fade out for 1 second

// This is equal to...
segment.FadeIn(0.5)
segment.FadeOut(1)
```

## Invert phase

```go
segment.InvertPhase()
```

## Reverse

```go
segment.Reverse()
```

## Repeat segment n times

```go
segment.Repeat(5)
```

## Change speed (speedup, slowdown)

```go
// Change speed by two times
segment.SetSpeed(2)

// Slow down
segment.SetSpeed(0.25)
```

## Take a slice out of a segment

```go
// Take a slice from the 1st second to the 2nd second (the length of the slice will be 1 second)
slice := segment.Slice(1, 2)
```

## Trim (cut)

```go
// Remove everything between the first second and the 1.5th second
segment.Trim(1, 1.5)

// Remove first 3 seconds
segment.TrimStart(3)

// Remove last 2.3 seconds
segment.TrimEnd(2.3)
```

## Empty segment

```go
// Sample rate, channels
gaudio.Empty(44100, 2)
```

## Silent segment

```go
// Returns a silent segment with the length of 5.5 seconds, 44100 sample rate, and two channels
gaudio.Silent(5.5, 44100, 2)
```

## Raw data

Raw data is stored in float32 in segment.Data, and is structured like that:

`[channel 1 value, channel 2 value, channel 1 value, channel 2 value, ...]`

## Audio formats

```go
const (
	FormatFLAC AudioFormat = 1 // .flac
	FormatWAVE AudioFormat = 2 // .wav
	FormatMP3  AudioFormat = 3 // .mp3
	FormatRAW  AudioFormat = 4 // Custom format

	// Tracker formats
	FormatMOD          AudioFormat = 5 // ProTracker .mod, rendered with gomodplay. This is faster than gotracker, but less accurate.
	FormatGotrackerMod AudioFormat = 6 // ProTracker .mod, rendered with Gotracker. This is slower than gomodplay, but more accurate.
	FormatS3M          AudioFormat = 7 // ScreamTracker III .s3m, rendered with Gotracker
	FormatXM           AudioFormat = 8 // FastTracker II .xm, rendered with Gotracker
	FormatIT           AudioFormat = 9 // ImpulseTracker .it, rendered with Gotracker
)
```

# TODO

-   Tests, benchmarks, profiling
-   Fix speed being rounded up to integers
-   More effects
-   Support more formats (while staying 100% Go)
