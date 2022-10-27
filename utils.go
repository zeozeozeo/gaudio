package gaudio

import (
	"fmt"
	"math"
)

// Returns the index of the specified second in segment.Data
// If the index is bigger than the segment length, returns the last index.
func (segment *Segment) GetIndexAtTime(time float64) uint64 {
	if time <= 0 {
		return 0
	}

	dataLen := len(segment.Data)
	if dataLen < 1 {
		dataLen = 1
	}

	maxIdx := uint64(dataLen - 1)
	idx := uint64(float64(segment.SampleRate)*time) * uint64(segment.Channels)
	if idx > maxIdx {
		return maxIdx
	}

	return idx
}

// Convert the sample rate of the segment.
func ConvertSampleRate(segment *Segment, want uint64) ([]float32, error) {
	if want == segment.SampleRate {
		return segment.Data, nil
	}

	if want < segment.SampleRate {
		newData := []float32{}
		dataLen := len(segment.Data)
		skipEvery := int(math.Ceil(float64(segment.SampleRate) / float64(want)))

		for i := 0; i < dataLen; i++ {
			if i%skipEvery == 0 {
				continue
			}
			newData = append(newData, segment.Data[i])
		}

		return newData, nil
	} else {
		return nil, fmt.Errorf(
			"TODO: converting to bigger sample rates is not yet implemented (got %d and %d)",
			segment.SampleRate,
			want,
		)
	}
}

func (segment *Segment) ConvertChannels(channels int) {
	// FIXME: This does not work if the new channels value is smaller
	// than the current one.
	if channels == segment.Channels || channels < 1 || channels < segment.Channels {
		return
	}

	newDataLen := int(segment.SampleRate) * channels
	newData := make([]float32, newDataLen)
	prevChannelValue := float32(0)

	for i := 0; i < newDataLen; i++ {
		if i%channels > i%segment.Channels {
			newData[i] = prevChannelValue
		} else {
			prevChannelValue = segment.Data[i]
			newData[i] = prevChannelValue
		}
	}
	segment.Data = newData
}

// Updates .Length and .LengthSeconds of the segment.
// Call after modifying the length of segment.Data.
func (segment *Segment) UpdateLength() {
	dataLen := len(segment.Data)
	segment.Length = uint64(dataLen) / uint64(segment.Channels)
	segment.LengthSeconds = float64(segment.Length) / float64(segment.SampleRate)
}

// Returns a copy of this segment.
func (segment *Segment) Clone() *Segment {
	dataLen := len(segment.Data)
	newSegment := &Segment{
		Data:          make([]float32, dataLen),
		Channels:      segment.Channels,
		SampleRate:    segment.SampleRate,
		Length:        segment.Length,
		LengthSeconds: segment.LengthSeconds,
		Format:        segment.Format,
	}

	// Copy audio data
	for i := 0; i < dataLen; i++ {
		newSegment.Data[i] = segment.Data[i]
	}

	return newSegment
}

func (segment *Segment) SamplesToSeconds(index uint64) float64 {
	return float64(index) / float64(segment.SampleRate)
}
