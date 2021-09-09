package wav

import (
	"bytes"
	"encoding/binary"
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

func (c WAVChunk) ReadHeaders(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, c.GetID().GetBytes())
	binary.Write(buffer, binary.LittleEndian, c.ChunkSize)
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

func (f WAVFormat) Encode(audio audio.Audio) []byte {

	buffer := new(bytes.Buffer)

	RIFFChunk{
		WAVChunk: WAVChunk{
			Chunk: utils.Chunk{
				ChunkID: RIFFID,
			},
			ChunkSize: 36 + uint32(len(audio.Samples)),
		},
		Format: WAVEID,
		FormatChunk: FormatChunk{
			WAVChunk: WAVChunk{
				Chunk: utils.Chunk{
					ChunkID: FORMATID,
				},
				ChunkSize: 16,
			},
			AudioFormat:   1,
			NumChannels:   audio.NumChannels,
			SampleRate:    uint32(audio.SampleRate),
			ByteRate:      uint32(audio.SampleRate) * uint32(audio.NumChannels) * uint32(audio.BitDepth) / 8,
			BlockAlign:    audio.NumChannels * audio.BitDepth / 8,
			BitsPerSample: audio.BitDepth,
		},
		DataChunk: DataChunk{
			WAVChunk: WAVChunk{
				Chunk: utils.Chunk{
					ChunkID: DATAID,
				},
				ChunkSize: uint32(len(audio.Samples)),
			},
			Data: audio.Samples,
		},
	}.Write(buffer)

	return buffer.Bytes()
}

func (f WAVFormat) IsFormat(data []byte) bool {

	// make sure headers match RIFF WAVE format
	_, _, _, err := RIFFHeaders(bytes.NewBuffer(data))
	return err == nil

}
