package aiff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var COMMONID utils.FourCC = []byte("COMM")

// TODO float80
type extended []byte

type CommonChunk struct {
	utils.Chunk
	NumChannels     int16
	NumSampleFrames uint32
	SampleSize      int16
	SampleRate      extended
}

func (c CommonChunk) GetBytes() []byte {
	return c.GetBytesWithHeaders(c.NumChannels, c.NumSampleFrames, c.SampleSize, c.SampleRate)
}

func NewCommonChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	// define common chunk struct
	var commChunk CommonChunk
	commChunk.ChunkID = COMMONID
	commChunk.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))

	// make sure common chunk size is 18
	if commChunk.ChunkSize != 18 {
		return CommonChunk{}, errors.New(fmt.Sprintf("COMMON chunk size is invalid: found %d, must be %d", commChunk.ChunkSize, 18))
	}

	// fill common chunk struct
	commChunk.NumChannels = int16(binary.BigEndian.Uint16(buffer.Next(2)))
	commChunk.NumSampleFrames = binary.BigEndian.Uint32(buffer.Next(4))
	commChunk.SampleSize = int16(binary.BigEndian.Uint16(buffer.Next(2)))
	commChunk.SampleRate = buffer.Next(10)

	return commChunk, nil
}
