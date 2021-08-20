package aiff

import (
	"bytes"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

type AIFFFormat struct {
	FormChunk FormChunk
}

type AIFFChunk struct {
	utils.Chunk
	ChunkSize int32
	// TODO data must be zero padded if odd
}

func (c AIFFChunk) GetBytesWithHeaders(fields ...interface{}) []byte {
	return c.GetBytesWithID(c.ChunkSize, fields)
}

func NewAIFFFormat(buffer *bytes.Buffer) (AIFFFormat, error) {

	// create form chunk
	formChunk, err := NewFormChunk(buffer)
	if err != nil {
		return AIFFFormat{}, err
	}

	return AIFFFormat{FormChunk: formChunk}, nil
}

// TODO multichannel table

// TODO optional chunks

func (f AIFFFormat) Decode(data []byte) (audio.Audio, error) {

	// create new AIFF format
	aiffFormat, err := NewAIFFFormat(bytes.NewBuffer(data))
	if err != nil {
		return audio.Audio{}, err
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
