package gaudio

// Detect silence (volume <= threshold) in audio.
// Returns the amount of silence from the start
// and the end (in seconds). Volume ranges from 0 to 1.
func (segment *Segment) DetectSilence(threshold float32) (float64, float64) {
	return segment.DetectSilenceStart(threshold), segment.DetectSilenceEnd(threshold)
}

// Detect starting silence (volume <= threshold) in audio in seconds
func (segment *Segment) DetectSilenceStart(threshold float32) float64 {
	dataLen := len(segment.Data)
	startSamples := 0

	// Detect silence from the start
	for i := 0; i < dataLen; i++ {
		absolute := segment.Data[i]
		if absolute < 0 {
			absolute = -absolute
		}

		if absolute <= threshold {
			startSamples++
		} else {
			break
		}
	}

	return segment.SamplesToSeconds(uint64(startSamples / segment.Channels))
}

// Detect trailing silence (volume <= threshold) in seconds
func (segment *Segment) DetectSilenceEnd(threshold float32) float64 {
	dataLen := len(segment.Data)
	endSamples := 0

	// Detect silence from the end
	for i := dataLen - 1; i >= 0; i-- {
		absolute := segment.Data[i]
		if absolute < 0 {
			absolute = -absolute
		}

		if absolute <= threshold {
			endSamples++
		} else {
			break
		}
	}

	return segment.SamplesToSeconds(uint64(endSamples / segment.Channels))
}

// Removes the starting silence (volume <= threshold) from the segment
func (segment *Segment) RemoveStartSilence(threshold float32) {
	segment.TrimStart(segment.DetectSilenceStart(threshold))
}

// Removes the trailing silence (volume <= threshold) from the segment
func (segment *Segment) RemoveEndSilence(threshold float32) {
	segment.TrimEnd(segment.DetectSilenceEnd(threshold))
}

// Removes all silence (volume <= threshold) from the segment
func (segment *Segment) RemoveStartAndEndSilence(threshold float32) {
	segment.RemoveStartSilence(threshold)
	segment.RemoveEndSilence(threshold)
}
