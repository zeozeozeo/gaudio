package gaudio

// Takes a slice out of this segment, and returns a new segment.
// start: the start of the slice in seconds
// end: the end of the slice in seconds
func (segment *Segment) Slice(start float64, end float64) *Segment {
	if start > end {
		start = end
	}
	if end < start {
		end = start
	}

	dataLen := uint64(len(segment.Data))
	startSample := uint64(start*float64(segment.SampleRate)) * uint64(segment.Channels)
	endSample := uint64(end*float64(segment.SampleRate)) * uint64(segment.Channels)
	lengthSamples := endSample - startSample
	// Make sure to keep the channels
	startSample -= uint64(segment.Channels % int(startSample))
	endSample -= uint64(segment.Channels % int(endSample))

	slicedAudio := make([]float32, lengthSamples)
	for i := startSample; i < dataLen && i < endSample; i++ {
		slicedAudio[i-startSample] = segment.Data[i]
	}

	newSegment := &Segment{
		Data:       slicedAudio,
		Channels:   segment.Channels,
		SampleRate: segment.SampleRate,
		Format:     segment.Format,
	}
	newSegment.UpdateLength()
	return newSegment
}
