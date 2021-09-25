package codec

import (
	"bytes"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/format"
)

func DecodeSpecific(data []byte, identifier format.FormatIdentifier) (a audio.Audio, e error) {

	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("panic while decoding format: %s", r)
		}
	}()

	aud, err := format.IdentifierToFormat[identifier].Decode(bytes.NewBuffer(data))
	if err != nil {
		return audio.Audio{}, fmt.Errorf("decode error: %s", err.Error())
	}

	return aud, nil

}

func Decode(data []byte) (audio.Audio, error) {

	// find matching format and decode
	for identifier, format := range format.IdentifierToFormat {
		if format.IsFormat(data) {
			return DecodeSpecific(data, identifier)
		}
	}

	return audio.Audio{}, fmt.Errorf("audio is either corrupted or is not supported")
}

func Encode(audio audio.Audio, identifier format.FormatIdentifier) []byte {
	buffer := new(bytes.Buffer)
	format.IdentifierToFormat[identifier].Encode(audio, buffer)
	return buffer.Bytes()
}
