package format

import (
	"bytes"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/format/aiff"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/format/wav"
)

type Format interface {
	Decode(data *bytes.Buffer) (audio.Audio, error)
	Encode(audio audio.Audio, buffer *bytes.Buffer)
	IsFormat(data []byte) bool
}

type FormatIdentifier int

const (
	WAV  FormatIdentifier = iota
	AIFF FormatIdentifier = iota
)

var IdentifierToFormat = map[FormatIdentifier]Format{
	WAV:  wav.WAVFormat{},
	AIFF: aiff.AIFFFormat{},
}
