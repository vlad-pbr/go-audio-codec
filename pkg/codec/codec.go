package codec

import (
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/format"
)

func Decode(data []byte, identifier format.FormatIdentifier) (audio.Audio, error) {
	return format.IdentifierToFormat[identifier].Decode(data)
}
