package wav

import (
	"bytes"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

type WAVFormat struct {
	RIFFChunk RIFFChunk
}

type WAVChunk struct {
	utils.Chunk
	ChunkSize uint32
}

func NewWAVFormat(buffer *bytes.Buffer) (WAVFormat, error) {

	// create riff chunk
	riffChunk, err := NewRIFFChunk(buffer)
	if err != nil {
		return WAVFormat{}, fmt.Errorf("error occurred while decoding RIFF chunk: %s", err.Error())
	}

	return WAVFormat{RIFFChunk: riffChunk}, nil
}

func (c WAVChunk) MakeChunkBytes(fields ...interface{}) []byte {
	return c.GetBytesWithID(
		c.ChunkSize,
		utils.GetBytes(false, fields),
	)
}

func (f WAVFormat) Decode(data []byte) (audio.Audio, error) {

	// create new WAVE format
	waveFormat, err := NewWAVFormat(bytes.NewBuffer(data))
	if err != nil {
		return audio.Audio{}, fmt.Errorf("error occurred while decoding WAV: %s", err.Error())
	}

	// fill audio struct with metadata
	audio := audio.Audio{
		NumChannels: waveFormat.RIFFChunk.FormatChunk.NumChannels,
		BitDepth:    waveFormat.RIFFChunk.FormatChunk.BitsPerSample,
		SampleRate:  uint64(waveFormat.RIFFChunk.FormatChunk.SampleRate),
		Samples:     waveFormat.RIFFChunk.DataChunk.Data,
	}

	return audio, nil
}

// TODO implement
func (f WAVFormat) Encode(audio audio.Audio) []byte {
	return []byte("")
}

func (f WAVFormat) IsFormat(data []byte) bool {

	// make sure headers match RIFF WAVE format
	_, _, _, err := RIFFHeaders(bytes.NewBuffer(data))
	return err == nil

}
