package aiff

import (
	"bytes"
	"encoding/binary"
	"fmt"

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

func NewAudioRecordingChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	var arChunk AudioRecordingChunk
	arChunk.ChunkID = AUDIORECORDINGID

	arChunk.ChunkSize = int32(binary.BigEndian.Uint32(utils.Next(buffer, 4)))
	if arChunk.ChunkSize != 24 {
		return arChunk, fmt.Errorf("%s chunk size is invalid: found %d, must be %d", string(AUDIORECORDINGID[:]), arChunk.ChunkSize, 24)
	}

	copy(arChunk.AESChannelStatusData[:], utils.Next(buffer, 24))

	return AudioRecordingChunk{}, nil
}
