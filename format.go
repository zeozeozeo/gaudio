package gaudio

type AudioFormat int

const (
	FormatFLAC AudioFormat = 1 // .flac
	FormatWAVE AudioFormat = 2 // .wav
	FormatMP3  AudioFormat = 3 // .mp3
	FormatRAW  AudioFormat = 4 // Custom format

	// Tracker formats
	FormatMOD          AudioFormat = 5 // ProTracker .mod, rendered with gomodplay. This is faster than gotracker, but less accurate.
	FormatGotrackerMod AudioFormat = 6 // ProTracker .mod, rendered with Gotracker. This is slower than gomodplay, but more accurate.
	FormatS3M          AudioFormat = 7 // ScreamTracker III .s3m, rendered with Gotracker
	FormatXM           AudioFormat = 8 // FastTracker II .xm, rendered with Gotracker
	FormatIT           AudioFormat = 9 // ImpulseTracker .it, rendered with Gotracker
)

func GetFormatString(format AudioFormat) string {
	switch format {
	case FormatWAVE:
		return "wav"
	case FormatFLAC:
		return "flac"
	case FormatMP3:
		return "mp3"
	case FormatRAW:
		return "raw"
	case FormatMOD, FormatGotrackerMod:
		return "mod"
	case FormatS3M:
		return "s3m"
	case FormatXM:
		return "xm"
	case FormatIT:
		return "it"
	}
	return ""
}
