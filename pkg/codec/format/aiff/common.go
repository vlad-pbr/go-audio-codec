package aiff

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils/float80"
)

var COMMONID utils.FourCC = [4]byte{67, 79, 77, 77} // COMM

type CommonChunk struct {
	AIFFChunk
	NumChannels     int16           // amount of audio channels
	NumSampleFrames uint32          // sample frame consists of sample per numChannels (= amount of samples / numChannels)
	SampleSize      int16           // NUMBER OF BITS for single audio sample (value can range from 1 to 32)
	SampleRate      float80.Float80 // SAMPLE FRAME (not samples themselves) played back / sec
}

func (c CommonChunk) GetBytes() []byte {
	return c.MakeChunkBytes(
		c.NumChannels,
		c.NumSampleFrames,
		c.SampleSize,
		c.SampleRate,
	)
}

func NewCommonChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	// define common chunk struct
	var commChunk CommonChunk
	commChunk.ChunkID = COMMONID
	commChunk.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))

	// make sure common chunk size is 18
	if commChunk.ChunkSize != 18 {
		return commChunk, fmt.Errorf("%s chunk size is invalid: found %d, must be %d", string(COMMONID[:]), commChunk.ChunkSize, 18)
	}

	// fill common chunk struct
	commChunk.NumChannels = int16(binary.BigEndian.Uint16(buffer.Next(2)))
	commChunk.NumSampleFrames = binary.BigEndian.Uint32(buffer.Next(4))
	commChunk.SampleSize = int16(binary.BigEndian.Uint16(buffer.Next(2)))

	// fill sample rate
	var sampleRateBytes [10]byte
	copy(sampleRateBytes[:], buffer.Next(10))
	commChunk.SampleRate = float80.NewFromBytes(sampleRateBytes)

	return commChunk, nil
}
