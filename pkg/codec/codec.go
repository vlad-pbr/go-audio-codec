package codec

import (
	"bytes"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/format"
)

func Decode(data []byte, identifier format.FormatIdentifier) (a audio.Audio, e error) {

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

func Encode(audio audio.Audio, identifier format.FormatIdentifier) []byte {
	buffer := new(bytes.Buffer)
	format.IdentifierToFormat[identifier].Encode(audio, buffer)
	return buffer.Bytes()
}
