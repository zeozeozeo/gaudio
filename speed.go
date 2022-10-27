package gaudio

import (
	"math"
)

// Speeds up the segment by n times without changing the sample rate.
// Currently this rounds up speeds to integers (read fixme).
func (segment *Segment) SpeedUp(speed float64) {
	if speed <= 1 {
		return
	}

	wantSampleRate := float64(segment.SampleRate) * speed

	// FIXME: This rounds up values to integers, so the speed will be rounded
	// 	      to an integer each time. Find a better way to handle this.
	skipEvery := uint64(math.Ceil(wantSampleRate / float64(segment.SampleRate)))

	newData := []float32{}
	dataLen := uint64(len(segment.Data))
	for i := uint64(0); i < dataLen-1; i++ {
		if i%skipEvery != 0 {
			continue
		}

		channel := i % uint64(segment.Channels)
		newData = append(newData, segment.Data[i+channel])
	}

	segment.Data = newData
	segment.UpdateLength()
}

// Slows down the segment by changing the sample rate.
func (segment *Segment) SlowDown(speed float64) {
	if speed >= 1 || speed <= 0 {
		return
	}

	segment.SampleRate = uint64(float64(segment.SampleRate) * speed)
	segment.UpdateLength()
}

func (segment *Segment) SetSpeed(speed float64) {
	if speed > 1 {
		// This won't change the sample rate
		segment.SpeedUp(speed)
	} else {
		// This will change the sample rate
		segment.SlowDown(speed)
	}
}
