package aiff

import (
	"bytes"
	"encoding/binary"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var SOUNDID utils.FourCC = [4]byte{83, 83, 78, 68} // SSND

type SoundDataChunk struct {
	AIFFChunk
	Offset    uint32
	BlockSize uint32
	SoundData []byte // sample frame size is always a multiple of 8
}

func (c SoundDataChunk) GetBytes() []byte {
	return c.MakeChunkBytes(
		c.Offset,
		c.BlockSize,
		c.SoundData,
	)
}

func NewSoundChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	// fill common chunk struct
	var soundChunk SoundDataChunk
	soundChunk.ChunkID = SOUNDID
	soundChunk.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))
	soundChunk.Offset = binary.BigEndian.Uint32(buffer.Next(4))
	soundChunk.BlockSize = binary.BigEndian.Uint32(buffer.Next(4))

	// read samples from buffer
	// actual semantics of these samples are only relevant when decoding to audio struct
	soundChunk.SoundData = buffer.Next(int(soundChunk.ChunkSize) - 8)

	AdjustForZeroPadding(soundChunk.ChunkSize, buffer)

	return soundChunk, nil
}
