package aiff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var SOUNDID utils.FourCC = []byte("SSND")

type SoundDataChunk struct {
	AIFFChunk
	Offset    uint32
	BlockSize uint32
	SoundData []uint8
}

func (c SoundDataChunk) GetBytes() []byte {
	return c.MakeChunkBytes(c.Offset, c.BlockSize, c.SoundData)
}

func NewSoundChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	// fill common chunk struct
	var soundChunk SoundDataChunk
	soundChunk.ChunkID = SOUNDID
	soundChunk.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))
	soundChunk.Offset = binary.BigEndian.Uint32(buffer.Next(4))
	soundChunk.BlockSize = binary.BigEndian.Uint32(buffer.Next(4))

	// parse sound chunk samples
	for i := 8; i != int(soundChunk.ChunkSize); i++ {
		sample, err := buffer.ReadByte()
		if err != nil {
			return SoundDataChunk{}, errors.New(fmt.Sprintf("unexpected EOF while reading SOUND chunk samples"))
		}
		soundChunk.SoundData = append(soundChunk.SoundData, uint8(sample))
	}

	AdjustForZeroPadding(soundChunk.ChunkSize, buffer)

	return soundChunk, nil
}
