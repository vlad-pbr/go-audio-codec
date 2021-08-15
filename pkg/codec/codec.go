package codec

import (
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/formats/wav"
)

type formatIdentifier int

const (
	WAV formatIdentifier = iota
)

type Format interface {
	GetName() string
}

var identifierToFormat = map[formatIdentifier]Format{
	WAV: wav.WAVFormat{},
}
