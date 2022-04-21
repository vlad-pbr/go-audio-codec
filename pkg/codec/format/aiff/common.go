package aiff

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils/extended"
)

var COMMONID utils.FourCC = [4]byte{67, 79, 77, 77} // COMM

type CommonChunk struct {
	AIFFChunk
	NumChannels     int16             // amount of audio channels
	NumSampleFrames uint32            // sample frame consists of sample per numChannels (= amount of samples / numChannels)
	SampleSize      int16             // NUMBER OF BITS for single audio sample (value can range from 1 to 32)
	SampleRate      extended.Extended // SAMPLE FRAME (not samples themselves) played back / sec
}

func (c CommonChunk) Write(buffer *bytes.Buffer) {
	c.WriteHeaders(buffer)
	binary.Write(buffer, binary.BigEndian, c.NumChannels)
	binary.Write(buffer, binary.BigEndian, c.NumSampleFrames)
	binary.Write(buffer, binary.BigEndian, c.SampleSize)
	binary.Write(buffer, binary.BigEndian, c.SampleRate.Bytes())
}

func NewCommonChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	// define common chunk struct
	var commChunk CommonChunk
	commChunk.ChunkID = COMMONID
	commChunk.ChunkSize = int32(binary.BigEndian.Uint32(utils.Next(buffer, 4)))

	// make sure common chunk size is 18
	if commChunk.ChunkSize != 18 {
		return commChunk, fmt.Errorf("%s chunk size is invalid: found %d, must be %d", string(COMMONID[:]), commChunk.ChunkSize, 18)
	}

	// fill common chunk struct
	commChunk.NumChannels = int16(binary.BigEndian.Uint16(utils.Next(buffer, 2)))
	commChunk.NumSampleFrames = binary.BigEndian.Uint32(utils.Next(buffer, 4))
	commChunk.SampleSize = int16(binary.BigEndian.Uint16(utils.Next(buffer, 2)))

	// fill sample rate
	var sampleRateBytes [10]byte
	copy(sampleRateBytes[:], utils.Next(buffer, 10))
	commChunk.SampleRate = extended.NewFromBytes(sampleRateBytes)

	return commChunk, nil
}
