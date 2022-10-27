package gaudio

// Overlay one segment on top of another.
// segment2 *gaudio.Segment: The segment you want to overlay.
// start float64: The second you want to overlay the segment at.
// This will automatically normalize the audio to avoid clipping.
// If the length of the overlayed segment is bigger, the second
// second will be cropped to the size of the target segment.
func (segment *Segment) Overlay(segment2 *Segment, start float64) error {
	// TODO: Make this work for negative values
	startIdx := segment.GetIndexAtTime(start)
	segmentDataLen := uint64(len(segment.Data))
	segment2DataLen := uint64(len(segment2.Data))

	// Convert sample rates if needed
	var convertedData []float32
	didConvert := false
	if segment.SampleRate != segment2.SampleRate {
		var err error
		convertedData, err = ConvertSampleRate(segment2, segment.SampleRate)
		if err != nil {
			return err
		}
		didConvert = true
		segment2DataLen = uint64(len(convertedData))
	}

	// Add values of two segments
	for i := startIdx; i < segmentDataLen && i-startIdx < segment2DataLen; i++ {
		if !didConvert {
			segment.Data[i] = segment.Data[i]/2 + segment2.Data[i-startIdx]/2
		} else {
			segment.Data[i] = segment.Data[i]/2 + convertedData[i-startIdx]/2
		}
	}

	// segment.Normalize()
	return nil
}
