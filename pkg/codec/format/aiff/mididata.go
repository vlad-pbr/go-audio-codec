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

func NewMIDIDataChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	// define chunk struct
	var midiChunk MIDIDataChunk
	midiChunk.ChunkID = MIDIDATAID
	midiChunk.ChunkSize = int32(binary.BigEndian.Uint32(utils.Next(buffer, 4)))
	midiChunk.MIDIData = utils.Next(buffer, int(midiChunk.ChunkSize))

	adjustForZeroPadding(midiChunk.ChunkSize, buffer)

	return midiChunk, nil

}
