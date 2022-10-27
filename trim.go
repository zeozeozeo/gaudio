package gaudio

// Removes everything between start..end from this segment (in seconds)
func (segment *Segment) Trim(start float64, end float64) {
	if start > end {
		start = end
	}
	trimStart := segment.GetIndexAtTime(start)
	trimEnd := segment.GetIndexAtTime(end)

	segment.Data = append(segment.Data[:trimStart], segment.Data[trimEnd:]...)
	segment.UpdateLength()
}

// Removes n seconds from the start of the segment
func (segment *Segment) TrimStart(seconds float64) {
	segment.Trim(0, seconds)
}

// Removs n seconds from the end of the segment
func (segment *Segment) TrimEnd(seconds float64) {
	segment.Trim(segment.LengthSeconds-seconds, segment.LengthSeconds)
}
