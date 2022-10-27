package gaudio

import (
	"fmt"
)

// Concatenate a segment at a specified index
func (segment *Segment) ConcatenateAt(segment2 *Segment, index uint64) error {
	if segment.Channels != segment2.Channels {
		return fmt.Errorf(
			"segments should have the same amount of channels"+
				" (got %d and %d). try using segment.ConvertChannels",
			segment.Channels,
			segment2.Channels,
		)
	}

	dataLen := uint64(len(segment.Data))
	if index > dataLen {
		index = dataLen
	}

	if segment.SampleRate == segment2.SampleRate {
		segment.Data = append(segment.Data[:index+1], append(segment2.Data, segment.Data[index:]...)...)
	} else {
		convertedData, err := ConvertSampleRate(segment2, segment.SampleRate)
		if err != nil {
			return err
		}
		segment.Data = append(segment.Data[:index+1], append(convertedData, segment.Data[index:]...)...)
	}

	segment.UpdateLength()
	segment.Normalize()
	return nil
}

func (segment *Segment) ConcatenateAtSecond(segment2 *Segment, second float64) {
	index := segment.GetIndexAtTime(second)
	segment.ConcatenateAt(segment2, index)
}

// Concatenate a segment at the end
func (segment *Segment) ConcatenateEnd(segment2 *Segment) error {
	return segment.ConcatenateAt(segment2, uint64(len(segment.Data)-1))
}

// Concatenate a segment at the start
func (segment *Segment) ConcatenateStart(segment2 *Segment) error {
	return segment.ConcatenateAt(segment2, 0)
}

// Repeats the segment n times
func (segment *Segment) Repeat(amount uint) error {
	oldSegment := *segment
	for i := uint(0); i < amount; i++ {
		err := segment.ConcatenateEnd(&oldSegment)
		if err != nil {
			return err
		}
	}
	return nil
}
