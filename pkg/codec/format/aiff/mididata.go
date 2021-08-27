package aiff

import (
	"bytes"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var MIDIDATAID utils.FourCC = [4]byte{77, 73, 68, 73} // MIDI

type MIDIDataChunk struct {
	AIFFChunk
	MIDIData []byte
}

func (c MIDIDataChunk) GetBytes() []byte {
	return c.MakeChunkBytes(
		c.MIDIData,
	)
}

// TODO implement
func NewMIDIDataChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return MIDIDataChunk{}, nil
}
