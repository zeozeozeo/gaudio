package gaudio

import (
	"math"
)

// Pitches audio down to the note.
// Pitch formula: 2 ** (key / 12).
// This will change the sample rate.
func (segment *Segment) ApplyNote(key int) {
	// TODO: This doesn't always work correctly
	// Convert key to pitch (2 ** (key / 12))
	pitch := math.Pow(2, math.Round(float64(key)/12))

	// Change sample rate
	segment.SampleRate = uint64(math.Round(float64(segment.SampleRate) * pitch))
}
