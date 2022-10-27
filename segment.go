package gaudio

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gotracker/gomixing/mixing"
	"github.com/gotracker/gomixing/sampling"
	tformat "github.com/gotracker/playback/format"
	"github.com/gotracker/playback/output"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"

	"github.com/hajimehoshi/go-mp3"
	"github.com/mewkiz/flac"
	"github.com/mjibson/go-dsp/wav"
	"github.com/zeozeozeo/gomodplay/pkg/mod"
)

type Segment struct {
	Data          []float32
	Channels      int     // Amount of channels (will be 2, even if there is one)
	SampleRate    uint64  // The sample rate of the segment
	Length        uint64  // Length of the segment in samples
	LengthSeconds float64 // Length of the segment in seconds
	Format        AudioFormat
}

var (
	errInvalidChecksum error = errors.New("invalid checksum")
)

// Loads an audio file from a reader.
func LoadAudio(reader io.Reader, format AudioFormat) (*Segment, error) {
	segment := &Segment{}

	switch format {
	case FormatMP3:
		err := decodeMp3(segment, reader)
		if err != nil {
			return nil, err
		}
	case FormatWAVE:
		err := decodeWav(segment, reader)
		if err != nil {
			return nil, err
		}
	case FormatFLAC:
		err := decodeFlac(segment, reader)
		if err != nil {
			return nil, err
		}
	case FormatMOD:
		err := decodeMod(segment, reader)
		if err != nil {
			return nil, err
		}
	case FormatGotrackerMod, FormatS3M, FormatXM, FormatIT:
		err := decodeTrackerModule(segment, reader, format)
		if err != nil {
			return nil, err
		}
	}
	// segment.Normalize()

	return segment, nil
}

// Loads an audio file from the file path.
func LoadAudioFromPath(path string, format AudioFormat) (*Segment, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadAudio(file, format)
}

// Converts a float32 array to a Segment.
// Structure of the audio array values (assume you have 2 channels):
// [channel 1 val, channel 2 val, channel 1 val, channel 2 val, ...]
func SegmentFromData(data []float32, channels int, sampleRate uint64) *Segment {
	if channels < 1 {
		channels = 1
	}
	if sampleRate < 1 {
		sampleRate = 1
	}

	segmentLength := uint64(len(data)) / uint64(channels)
	segment := &Segment{
		Data:          data,
		Channels:      channels,
		SampleRate:    sampleRate,
		Length:        segmentLength,
		LengthSeconds: float64(segmentLength) / float64(sampleRate),
		Format:        FormatRAW,
	}
	// segment.Normalize()

	return segment
}

// Returns a silent segment, with the length of n seconds.
func Silent(seconds float64, sampleRate uint64, channels uint) *Segment {
	if sampleRate < 1 {
		sampleRate = 1
	}
	if seconds < 0 {
		seconds = 0
	}

	dataLength := uint64(seconds*float64(sampleRate)) * uint64(channels)
	return &Segment{
		Data:          make([]float32, dataLength),
		Channels:      int(channels),
		SampleRate:    sampleRate,
		Length:        dataLength / uint64(channels),
		LengthSeconds: (float64(dataLength) / float64(channels)) / float64(sampleRate),
		Format:        FormatRAW,
	}
}

// Returns an empty segment.
func Empty(sampleRate uint64, channels uint) *Segment {
	return &Segment{
		Channels:      int(channels),
		SampleRate:    sampleRate,
		Length:        0,
		LengthSeconds: 0,
		Format:        FormatRAW,
	}
}

func convertBytesToFloat32(data []byte, dataLength uint64) []float32 {
	resultLength := dataLength / uint64(2)
	result := make([]float32, resultLength)

	for i := uint64(0); i < resultLength; i++ {
		s16 := int16(binary.LittleEndian.Uint16(data[2*i : 2*i+2]))
		float := float32(s16) / 0x7FFF
		result[i] = float
	}

	return result
}

func decodeMp3(segment *Segment, reader io.Reader) error {
	decoder, err := mp3.NewDecoder(reader)
	if err != nil {
		return err
	}

	segment.SampleRate = uint64(decoder.SampleRate())
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(decoder)
	if err != nil {
		return err
	}

	dataLength := uint64(decoder.Length())
	segment.Data = convertBytesToFloat32(buf.Bytes(), dataLength)
	segment.Channels = 2 // go-mp3 always outputs 2 channels
	segment.Length = dataLength / 4
	segment.LengthSeconds = float64(segment.Length) / float64(segment.SampleRate)
	segment.Format = FormatMP3

	return nil
}

func decodeWav(segment *Segment, reader io.Reader) error {
	decoder, err := wav.New(reader)
	if err != nil {
		return err
	}

	segment.Data, err = decoder.ReadFloats(decoder.Samples)
	if err != nil {
		fmt.Printf("died %s", err)
		return err
	}

	segment.SampleRate = uint64(decoder.SampleRate)
	segment.Channels = int(decoder.NumChannels)
	segment.Length = uint64(decoder.Samples / segment.Channels)
	segment.LengthSeconds = float64(segment.Length) / float64(segment.SampleRate)
	segment.Format = FormatWAVE

	return nil
}

func decodeFlac(segment *Segment, reader io.Reader) error {
	stream, err := flac.New(reader)
	if err != nil {
		return err
	}

	segment.Channels = int(stream.Info.NChannels)
	segment.SampleRate = uint64(stream.Info.SampleRate)
	segment.Length = stream.Info.NSamples
	segment.LengthSeconds = float64(segment.Length) / float64(segment.SampleRate)
	segment.Format = FormatFLAC

	// TODO: Find a better way to do this
	dataLength := segment.Length * uint64(segment.Channels)
	channelsData := make([][]float32, segment.Channels)
	segment.Data = make([]float32, dataLength)

	md5sum := md5.New()

	for {
		frame, err := stream.ParseNext()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		frame.Hash(md5sum)

		// Iterate through all channels (subframes)
		for channelNum, subframe := range frame.Subframes {
			for _, sample := range subframe.Samples {
				channelsData[channelNum] = append(channelsData[channelNum], float32(sample)/0x7FFF)
			}
		}
	}

	gotSum, wantSum := md5sum.Sum(nil), stream.Info.MD5sum[:]
	if !bytes.Equal(gotSum, wantSum) {
		return errInvalidChecksum
	}

	for i := uint64(0); i < dataLength; i++ {
		segment.Data[i] = channelsData[i%uint64(segment.Channels)][i/uint64(segment.Channels)]
	}

	return nil
}

func decodeMod(segment *Segment, reader io.Reader) error {
	sampleRate := 44100
	player := mod.NewModPlayer(uint32(sampleRate))
	err := player.LoadModFile(reader)
	if err != nil {
		return err
	}

	segment.Channels = 2
	segment.SampleRate = uint64(sampleRate)
	segment.Format = FormatMOD

	player.Play()
	for {
		samples := make([][2]float32, 1)
		_, ok := player.Stream(samples)
		if !ok || player.State.HasLooped {
			break
		}

		segment.Data = append(segment.Data, samples[0][0], samples[0][1])
		segment.Length++
	}

	segment.LengthSeconds = float64(segment.Length) / float64(sampleRate)
	segment.Normalize()
	return nil
}

// Decodes .xm, .it, .s3m, .mod modules (with Gotracker)
func decodeTrackerModule(segment *Segment, reader io.Reader, format AudioFormat) error {
	var features []feature.Feature

	features = append(features,
		feature.UseNativeSampleFormat(true),
		feature.IgnoreUnknownEffect{Enabled: true},
		feature.SongLoop{Count: 0},
	)

	// Make a readseeker
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	readseeker := FakeReadSeeker(buf.Bytes())

	// Setup player
	player, _, err := tformat.LoadFromReader(GetFormatString(format), &readseeker, features...)
	if err != nil {
		return err
	}

	if err := player.SetupSampler(44100, 2); err != nil {
		return err
	}

	if err := player.Configure(features); err != nil {
		return err
	}

	premixDataChannel := make(chan *output.PremixData, 8)
	defer close(premixDataChannel)

	// Create mixers
	mixer := mixing.Mixer{
		Channels: 2,
	}
	panMixer := mixing.GetPanMixer(2)

	go func() {
		// Wait for data to appear
		for premix := range premixDataChannel {
			data := mixer.Flatten(
				panMixer,
				premix.SamplesLen,
				premix.Data,
				premix.MixerVolume,
				sampling.Format16BitLESigned,
			)

			// Convert the data to float32 and append it to the segment
			segment.Data = append(segment.Data, convertBytesToFloat32(data, uint64(len(data)))...)
		}
	}()

	// Render the song
	for {
		if err := player.Update(0, premixDataChannel); err != nil {
			if errors.Is(err, song.ErrStopSong) {
				break
			}
			return err
		}
	}

	segment.Channels = 2
	segment.Format = format
	segment.SampleRate = 44100
	segment.UpdateLength()

	return nil
}
