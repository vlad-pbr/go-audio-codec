package aiff

import (
	"bytes"
	"encoding/binary"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var AUDIORECORDINGID utils.FourCC = [4]byte{65, 69, 83, 68} // AESD

type AudioRecordingChunk struct { // size is always 24
	AIFFChunk
	AESChannelStatusData [24]byte
}

func (c AudioRecordingChunk) Write(buffer *bytes.Buffer) {
	c.WriteHeaders(buffer)
	binary.Write(buffer, binary.BigEndian, c.AESChannelStatusData[:])
}

// TODO implement
func NewAudioRecordingChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return AudioRecordingChunk{}, nil
}
