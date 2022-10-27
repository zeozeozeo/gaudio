package gaudio

// Reverses this segment.
func (segment *Segment) Reverse() {
	dataLen := len(segment.Data)

	for i, j := 0, dataLen-1; i < j; i, j = i+1, j-1 {
		segment.Data[i], segment.Data[j] = segment.Data[j], segment.Data[i]
	}
}

// Inverts the phase of this segment
func (segment *Segment) InvertPhase() {
	dataLen := len(segment.Data)
	for i := 0; i < dataLen; i++ {
		segment.Data[i] *= -1
	}
}
