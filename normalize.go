package gaudio

// Normalizes audio values of this segment from -1.0 to 1.0.
func (segment *Segment) Normalize() {
	dataLen := len(segment.Data)
	maxValue := segment.GetMaxAbsValue()

	// Normalize
	for i := 0; i < dataLen; i++ {
		segment.Data[i] /= maxValue
	}
}

// Get the maximum (absolute) value of the segment.
func (segment *Segment) GetMaxAbsValue() float32 {
	maxValue := float32(0)
	for _, val := range segment.Data {
		if val < 0 {
			val *= -1
		}

		if val > maxValue {
			maxValue = val
		}
	}

	return maxValue
}
