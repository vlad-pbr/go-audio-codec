package aiff

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

type AIFFFormat struct {
}

type AIFFChunk struct {
	utils.Chunk
	ChunkSize int32 // in size is odd - zero pad must be present
}

func (c AIFFChunk) WriteHeaders(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, c.GetID().GetBytes())
	binary.Write(buffer, binary.BigEndian, c.ChunkSize)
}

func adjustForZeroPadding(size int32, buffer *bytes.Buffer, pad bool) {

	// drop/add zero pad byte if chunk size is odd
	if size%2 != 0 {
		if pad {
			buffer.Write([]byte{0})
		} else {
			utils.Next(buffer, 1)
		}
	}

}

func (f AIFFFormat) Decode(data *bytes.Buffer) (audio.Audio, error) {

	// create form chunk
	formChunk, err := NewFormChunk(data)
	if err != nil {
		return audio.Audio{}, fmt.Errorf("error occurred while decoding FORM chunk: %s", err.Error())
	}

	var commonChunkIndex int
	var soundChunkIndex int

	// find required form local chunks
	for index, chunk := range formChunk.LocalChunks {

		chunkID := chunk.GetID()

		switch string(chunkID[:]) {
		case string(COMMONID[:]):
			commonChunkIndex = index
		case string(SOUNDID[:]):
			soundChunkIndex = index
		}
	}

	// calculate samplerate from extended precision float bytes
	sampleRate, _ := formChunk.LocalChunks[commonChunkIndex].(CommonChunk).SampleRate.Float().Uint64()
	samplesLen := len(formChunk.LocalChunks[soundChunkIndex].(SoundDataChunk).SoundData)
	samplesOffset := int(formChunk.LocalChunks[soundChunkIndex].(SoundDataChunk).Offset)

	// generate audio container
	return audio.NewAudio(
		uint16(formChunk.LocalChunks[commonChunkIndex].(CommonChunk).NumChannels),
		sampleRate,
		uint16(formChunk.LocalChunks[commonChunkIndex].(CommonChunk).SampleSize),
		formChunk.LocalChunks[soundChunkIndex].(SoundDataChunk).SoundData[samplesOffset:samplesLen-samplesOffset],
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
