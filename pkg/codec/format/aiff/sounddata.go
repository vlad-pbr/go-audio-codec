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
	utils.Chunk
	offset    uint32
	blockSize uint32
	soundData []uint8
}

func (c SoundDataChunk) GetBytes() []byte {
	return c.GetBytesWithHeaders(c.offset, c.blockSize, c.soundData)
}

func NewSoundChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	// fill common chunk struct
	var soundChunk SoundDataChunk
	soundChunk.ChunkID = SOUNDID
	soundChunk.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))
	soundChunk.offset = binary.BigEndian.Uint32(buffer.Next(4))
	soundChunk.blockSize = binary.BigEndian.Uint32(buffer.Next(4))

	// parse sound chunk samples
	for i := 8; i != int(soundChunk.ChunkSize); i++ {
		sample, err := buffer.ReadByte()
		if err != nil {
			return SoundDataChunk{}, errors.New(fmt.Sprintf("unexpected EOF while reading SOUND chunk samples"))
		}
		soundChunk.soundData = append(soundChunk.soundData, uint8(sample))
	}

	return soundChunk, nil
}
