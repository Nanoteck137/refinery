package types

type MediaFormat string

const (
	MediaFormatUnknown  MediaFormat = "unknown"
	MediaFormatEmpty    MediaFormat = ""
	MediaFormatFlac     MediaFormat = "flac"
	MediaFormatPcmS16LE MediaFormat = "pcm_s16le"
	MediaFormatOpus     MediaFormat = "opus"
	MediaFormatVorbis   MediaFormat = "vorbis"
	MediaFormatMp3      MediaFormat = "mp3"
	MediaFormatAac      MediaFormat = "aac"
)

var ValidMediaFormats = []MediaFormat{
	MediaFormatFlac,
	MediaFormatPcmS16LE,
	MediaFormatOpus,
	MediaFormatVorbis,
	MediaFormatMp3,
	MediaFormatAac,
}

type MediaInfo struct {
	Name       string
	Ext        string
	IsLossless bool
	Order      int
}

var MediaFormatInfos = map[MediaFormat]MediaInfo{
	MediaFormatFlac: {
		Name:       "FLAC",
		Ext:        ".flac",
		IsLossless: true,
		Order:      0,
	},
	MediaFormatPcmS16LE: {
		Name:       "PCM-S16-LE",
		Ext:        ".wav",
		IsLossless: true,
		Order:      1,
	},
	MediaFormatOpus: {
		Name:  "Opus",
		Ext:   ".opus",
		Order: 2,
	},
	MediaFormatVorbis: {
		Name:  "Vorbis",
		Ext:   ".ogg",
		Order: 3,
	},
	MediaFormatMp3: {
		Name:  "MP3",
		Ext:   ".mp3",
		Order: 4,
	},
	MediaFormatAac: {
		Name:  "AAC",
		Ext:   ".m4a",
		Order: 5,
	},
}

func (m MediaFormat) ToExt() (string, bool) {
	switch m {
	case MediaFormatFlac:
		return ".flac", true
	case MediaFormatPcmS16LE:
		return ".wav", true
	case MediaFormatOpus:
		return ".opus", true
	case MediaFormatVorbis:
		return ".ogg", true
	case MediaFormatMp3:
		return ".mp3", true
	case MediaFormatAac:
		return ".m4a", true
	}

	return "", false
}

func (m MediaFormat) IsValid() bool {
	switch m {
	case MediaFormatFlac:
		return true
	case MediaFormatPcmS16LE:
		return true
	case MediaFormatOpus:
		return true
	case MediaFormatVorbis:
		return true
	case MediaFormatMp3:
		return true
	case MediaFormatAac:
		return true
	}

	return false
}

func (m MediaFormat) IsLossless() bool {
	switch m {
	case MediaFormatFlac:
		return true
	case MediaFormatPcmS16LE:
		return true
	}

	return false
}

func (m MediaFormat) IsLossy() bool {
	switch m {
	case MediaFormatOpus:
		return true
	case MediaFormatVorbis:
		return true
	case MediaFormatMp3:
		return true
	case MediaFormatAac:
		return true
	}

	return false
}
