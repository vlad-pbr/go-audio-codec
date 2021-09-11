package aiff

import (
	"bytes"
	"encoding/binary"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var MIDIDATAID utils.FourCC = [4]byte{77, 73, 68, 73} // MIDI

type MIDIDataChunk struct {
	AIFFChunk
	MIDIData []byte
}

func (c MIDIDataChunk) Write(buffer *bytes.Buffer) {
	c.WriteHeaders(buffer)
	binary.Write(buffer, binary.BigEndian, c.MIDIData)
}

// TODO implement
func NewMIDIDataChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return MIDIDataChunk{}, nil
}
