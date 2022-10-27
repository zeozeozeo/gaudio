package gaudio

import (
	"math"
	"math/rand"
)

const (
	bpmInterval = 128
)

// Calculates the BPM of this segment.
// Warning: this is very slow (and inaccurate), and even slower for bigger segments.
// Use a smaller slice of a segment instead of the entire thing to get faster
// results.
// minBpm: the minimum BPM you expect (120)
// maxBpm: the maximum BPM you expect (200)
// steps: the amount of steps per interval (1024)
// samplesPerBeat: the amount of samples used for a single beat (1024)
// Returns the calculated BPM value
func (segment *Segment) CalcBPM(minBpm float64, maxBpm float64, steps float64, samplesPerBeat uint64) float64 {
	slowest := bpmToInterval(minBpm, samplesPerBeat)
	fastest := bpmToInterval(maxBpm, samplesPerBeat)
	step := (slowest - fastest) / steps

	height := math.Inf(0)
	trough := math.NaN()
	dataLen := len(segment.Data)

	for interval := fastest; interval <= slowest; interval += step {
		total := float64(0)

		for s := 0; s < dataLen; s++ {
			total += intervalDifference(segment, interval)
		}

		if total < height {
			trough = interval
			height = total
		}
	}

	return intervalToBpm(trough, samplesPerBeat)
}

func bpmToInterval(bpm float64, sampleRate uint64) float64 {
	var beatsPerSecond, samplesPerBeat float64

	beatsPerSecond = bpm / 60
	samplesPerBeat = float64(sampleRate) / beatsPerSecond
	return samplesPerBeat / bpmInterval
}

func intervalDifference(segment *Segment, interval float64) float64 {
	var diff, total float64
	beats := [...]float64{-32, -16, -8, -4, -2, -1, 1, 2, 4, 8, 16, 32}
	nobeats := [...]float64{-0.5, -0.25, 0.25, 0.5}

	mid := rand.Float64() * float64(len(segment.Data))
	v := sample(segment, mid)

	diff, total = 0.0, 0.0

	for n := 0; n < (len(beats) / 2); n++ {
		y := sample(segment, mid+beats[n]*interval)
		w := 1.0 / math.Abs(beats[n])

		diff += w * math.Abs(y-v)
		total += w
	}

	for n := 0; n < (len(nobeats) / 2); n++ {
		y := sample(segment, mid+nobeats[n]*interval)
		w := math.Abs(nobeats[n])

		diff -= w * math.Abs(y-v)
		total += w
	}

	return diff / total
}

func sample(segment *Segment, offset float64) float64 {
	n := math.Floor(offset)
	i := int64(n)

	if n >= 0.0 && n < float64(len(segment.Data)) {
		return float64(segment.Data[i])
	}
	return 0.0
}

func intervalToBpm(interval float64, sampleRate uint64) float64 {
	var samplesPerBeat, beatsPerSecond float64

	samplesPerBeat = interval * bpmInterval
	beatsPerSecond = float64(sampleRate) / samplesPerBeat
	return beatsPerSecond * 60
}
