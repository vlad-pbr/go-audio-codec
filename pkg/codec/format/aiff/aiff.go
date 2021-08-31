package aiff

import (
	"bytes"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils/float80"
)

type AIFFFormat struct {
	FormChunk FormChunk
}

type AIFFChunk struct {
	utils.Chunk
	ChunkSize int32 // in size is odd - zero pad must be present
}

func (c AIFFChunk) MakeChunkBytes(fields ...interface{}) []byte {
	return c.GetBytesWithID(
		c.ChunkSize,
		utils.GetBytes(true, fields),
	)
}

func NewAIFFFormat(buffer *bytes.Buffer) (AIFFFormat, error) {

	// create form chunk
	formChunk, err := NewFormChunk(buffer)
	if err != nil {
		return AIFFFormat{}, fmt.Errorf("error occurred while decoding FORM chunk: %s", err.Error())
	}

	return AIFFFormat{FormChunk: formChunk}, nil
}

func AdjustForZeroPadding(size int32, buffer *bytes.Buffer) {

	// drop zero pad byte if chunk size is odd
	if size%2 != 0 {
		buffer.Next(1)
	}

}

// TODO multichannel table

// TODO optional chunks

func (f AIFFFormat) Decode(data []byte) (audio.Audio, error) {

	// create new AIFF format
	aiffFormat, err := NewAIFFFormat(bytes.NewBuffer(data))
	if err != nil {
		return audio.Audio{}, fmt.Errorf("error occurred while decoding AIFF: %s", err.Error())
	}

	var commonChunkIndex int
	var audioChunkIndex int

	// find required form local chunks
	for index, chunk := range aiffFormat.FormChunk.LocalChunks {

		chunkID := chunk.GetID()

		switch string(chunkID[:]) {
		case string(COMMONID[:]):
			commonChunkIndex = index
		case string(SOUNDID[:]):
			audioChunkIndex = index
			_ = audioChunkIndex // TODO stub
		}
	}

	// fill audio struct with metadata
	audio := audio.Audio{
		NumChannels: uint16(aiffFormat.FormChunk.LocalChunks[commonChunkIndex].(CommonChunk).NumChannels),
		BitDepth:    uint16(aiffFormat.FormChunk.LocalChunks[commonChunkIndex].(CommonChunk).SampleSize),
		SampleRate:  float80.NewFromBytes(aiffFormat.FormChunk.LocalChunks[commonChunkIndex].(CommonChunk).SampleRate).UInt64(),
	}

	// TODO read sample data in soundchunk

	return audio, nil
}

// TODO implement
func (f AIFFFormat) Encode(audio.Audio) []byte {
	return []byte("")
}
