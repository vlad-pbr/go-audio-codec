package aiff

import (
	"bytes"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
)

type AIFFFormat struct {
}

// TODO multichannel table

// TODO optional chunks

func (f AIFFFormat) Decode(data []byte) (audio.Audio, error) {

	// create form chunk
	formChunk, err := NewFormChunk(bytes.NewBuffer(data))
	if err != nil {
		return audio.Audio{}, err
	}

	// define audio struct
	audio := audio.Audio{}

	// iterate form local chunks and fill audio struct accordingly
	for _, chunk := range formChunk.localChunks {
		if bytes.Compare(chunk.GetID(), COMMONID) == 0 {
			audio.NumChannels = uint16(chunk.(CommonChunk).numChannels)
		}
	}

	return audio, nil
}
