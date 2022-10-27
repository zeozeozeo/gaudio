package gaudio

// Applies an echo effect on the segment.
// delay: delay between echo layers (0.05)
// layers: amount of layer (2)
func (segment *Segment) EchoEffect(delay float64, layers int) {
	lastDelay := float64(0)
	originalSegment := *segment

	for l := 0; l < int(layers); l++ {
		lastDelay += delay
		segment.Overlay(&originalSegment, lastDelay)
	}
}
