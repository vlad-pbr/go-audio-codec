package aiff

import (
	"bytes"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var AUDIORECORDINGID utils.FourCC = [4]byte{65, 69, 83, 68} // AESD

type AudioRecordingChunk struct { // size is always 24
	AIFFChunk
	AESChannelStatusData [24]byte
}

func (c AudioRecordingChunk) GetBytes() []byte {
	return c.MakeChunkBytes(
		c.AESChannelStatusData,
	)
}

// TODO implement
func NewAudioRecordingChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return AudioRecordingChunk{}, nil
}
