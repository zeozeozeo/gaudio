package gaudio

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"

	"github.com/youpy/go-wav"
)

var (
	errUnsupportedExportFormat error = errors.New("unsupported export format")
)

func (segment *Segment) Export(writer io.Writer, format AudioFormat) error {
	switch format {
	case FormatWAVE:
		segmentLength := len(segment.Data)
		wavWriter := wav.NewWriter(writer, uint32(segmentLength*2), uint16(segment.Channels), uint32(segment.SampleRate), 16)
		writer := bufio.NewWriter(wavWriter) // Use a buffered writer

		for _, sample := range segment.Data {
			// Clamp sample to -1 .. 1
			if sample > 1 {
				sample = 1
			} else if sample < -1 {
				sample = -1
			}

			// Convert float32 to signed int16
			var s16 int16
			if sample < 0 {
				s16 = int16(sample * 0x8000)
			} else {
				s16 = int16(sample * 0x7FFF)
			}

			err := binary.Write(writer, binary.LittleEndian, s16)
			if err != nil {
				return err
			}
		}
		writer.Flush()
	default:
		return errUnsupportedExportFormat
	}
	return nil
}
