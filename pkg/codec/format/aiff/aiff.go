package aiff

import (
	"bytes"
	"encoding/binary"
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

func (c AIFFChunk) WriteHeaders(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, c.GetID().GetBytes())
	binary.Write(buffer, binary.BigEndian, c.ChunkSize)
}

func NewAIFFFormat(buffer *bytes.Buffer) (AIFFFormat, error) {

	// create form chunk
	formChunk, err := NewFormChunk(buffer)
	if err != nil {
		return AIFFFormat{}, fmt.Errorf("error occurred while decoding FORM chunk: %s", err.Error())
	}

	return AIFFFormat{FormChunk: formChunk}, nil
}

func adjustForZeroPadding(size int32, buffer *bytes.Buffer) {

	// drop zero pad byte if chunk size is odd
	if size%2 != 0 {
		utils.Next(buffer, 1)
	}

}

func (f AIFFFormat) Decode(data *bytes.Buffer) (audio.Audio, error) {

	// create new AIFF format
	aiffFormat, err := NewAIFFFormat(data)
	if err != nil {
		return audio.Audio{}, fmt.Errorf("error occurred while decoding AIFF: %s", err.Error())
	}

	var commonChunkIndex int
	var soundChunkIndex int

	// find required form local chunks
	for index, chunk := range aiffFormat.FormChunk.LocalChunks {

		chunkID := chunk.GetID()

		switch string(chunkID[:]) {
		case string(COMMONID[:]):
			commonChunkIndex = index
		case string(SOUNDID[:]):
			soundChunkIndex = index
		}
	}

	// calculate samplerate from extended precision float bytes
	sampleRate, _ := aiffFormat.FormChunk.LocalChunks[commonChunkIndex].(CommonChunk).SampleRate.Float().Uint64()
	samplesLen := len(aiffFormat.FormChunk.LocalChunks[soundChunkIndex].(SoundDataChunk).SoundData)
	samplesOffset := int(aiffFormat.FormChunk.LocalChunks[soundChunkIndex].(SoundDataChunk).Offset)

	// generate audio container
	return audio.NewAudio(
		uint16(aiffFormat.FormChunk.LocalChunks[commonChunkIndex].(CommonChunk).NumChannels),
		sampleRate,
		uint16(aiffFormat.FormChunk.LocalChunks[commonChunkIndex].(CommonChunk).SampleSize),
		aiffFormat.FormChunk.LocalChunks[soundChunkIndex].(SoundDataChunk).SoundData[samplesOffset:samplesLen-samplesOffset],
		binary.BigEndian,
	)
}

// TODO implement
func (f AIFFFormat) Encode(audio audio.Audio, buffer *bytes.Buffer) {
}

func (f AIFFFormat) IsFormat(data []byte) bool {

	// make sure headers match FORM AIFF format
	_, _, _, err := FormHeaders(bytes.NewBuffer(data))
	return err == nil

}
