package gaudio

// Fade in/out effect.
// start float64: the amount of fade-in seconds
// end float64: the amount of fade-out seconds
func (segment *Segment) Fade(start float64, end float64) {
	fadeStart := segment.GetIndexAtTime(start)
	fadeEnd := segment.GetIndexAtTime(segment.LengthSeconds - end)
	dataLen := uint64(len(segment.Data))

	for i := uint64(0); i < fadeStart && i < dataLen; i++ {
		segment.Data[i] = segment.Data[i] * (float32(i) / float32(fadeStart))
	}

	fadeMul := float32(0)
	for i := dataLen - 1; i > fadeEnd; i-- {
		segment.Data[i] = segment.Data[i] * fadeMul
		fadeMul += 1 / float32(dataLen-fadeEnd)
	}
}

// Fade in effect. Use Segment.Fade if you want both fade in and out.
func (segment *Segment) FadeIn(seconds float64) {
	segment.Fade(seconds, 0)
}

// Fade out effect. Use Segment.Fade if you want both fade in and out.
func (segment *Segment) FadeOut(seconds float64) {
	segment.Fade(0, seconds)
}
