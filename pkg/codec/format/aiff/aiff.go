package aiff

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

type AIFFFormat struct {
	FormChunk FormChunk
}

type AIFFChunk struct {
	utils.Chunk
	ChunkSize int32 // in size is odd - zero pad must be present
}

func (c AIFFChunk) MakeChunkBytes(fields ...interface{}) []byte {
	return c.GetBytesWithID(c.ChunkSize, utils.GetBytes(true, fields))
}

func NewAIFFFormat(buffer *bytes.Buffer) (AIFFFormat, error) {

	// create form chunk
	formChunk, err := NewFormChunk(buffer)
	if err != nil {
		return AIFFFormat{}, errors.New(fmt.Sprintf("error occurred while decoding FORM chunk: %s", err.Error()))
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
		return audio.Audio{}, errors.New(fmt.Sprintf("error occurred while decoding AIFF: %s", err.Error()))
	}

	// define audio struct
	audio := audio.Audio{}

	// TODO iterate form local chunks and fill audio struct accordingly
	for _, chunk := range aiffFormat.FormChunk.LocalChunks {
		if bytes.Compare(chunk.GetID(), COMMONID) == 0 {
			audio.NumChannels = uint16(chunk.(CommonChunk).NumChannels)
		}
	}

	return audio, nil
}

// TODO implement
func (f AIFFFormat) Encode(audio.Audio) []byte {
	return []byte("")
}
