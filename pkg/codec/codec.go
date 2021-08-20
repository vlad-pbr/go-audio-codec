package codec

import (
	"errors"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/format"
)

func Decode(data []byte, identifier format.FormatIdentifier) (audio.Audio, error) {
	aud, err := format.IdentifierToFormat[identifier].Decode(data)
	if err != nil {
		return audio.Audio{}, errors.New(fmt.Sprintf("decode error: %s", err.Error()))
	}

	return aud, nil
}

func Encode(audio audio.Audio, identifier format.FormatIdentifier) []byte {
	return format.IdentifierToFormat[identifier].Encode(audio)
}
