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

func (c WAVChunk) WriteHeaders(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, c.GetID().GetBytes())
	binary.Write(buffer, binary.LittleEndian, c.ChunkSize)
}

func (f WAVFormat) Decode(data *bytes.Buffer) (audio.Audio, error) {

	// create new WAVE format
	waveFormat, err := NewWAVFormat(data)
	if err != nil {
		return audio.Audio{}, fmt.Errorf("error occurred while decoding WAV: %s", err.Error())
	}

	// audio container out of wave format
	return audio.NewAudio(
		waveFormat.RIFFChunk.FormatChunk.NumChannels,
		uint64(waveFormat.RIFFChunk.FormatChunk.SampleRate),
		waveFormat.RIFFChunk.FormatChunk.BitsPerSample,
		waveFormat.RIFFChunk.DataChunk.Data,
		binary.LittleEndian,
	)
}

func (f WAVFormat) Encode(audio audio.Audio, buffer *bytes.Buffer) {

	// round bit depth to nearest multiple of 8
	bitDepth := audio.BitDepth()
	if bitDepth%8 != 0 {
		bitDepth = bitDepth + (bitDepth % 8)
	}

	RIFFChunk{
		WAVChunk: WAVChunk{
			Chunk: utils.Chunk{
				ChunkID: RIFFID,
			},
			ChunkSize: 36 + uint32(audio.BytesAmount()),
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
			NumChannels:   audio.NumChannels(),
			SampleRate:    uint32(audio.SampleRate()),
			ByteRate:      uint32(audio.SampleRate()) * uint32(audio.NumChannels()) * uint32(audio.ByteDepth()),
			BlockAlign:    audio.NumChannels() * uint16(audio.ByteDepth()),
			BitsPerSample: bitDepth,
		},
		DataChunk: DataChunk{
			WAVChunk: WAVChunk{
				Chunk: utils.Chunk{
					ChunkID: DATAID,
				},
				ChunkSize: uint32(audio.BytesAmount()),
			},
			Data: audio.Samples(binary.LittleEndian),
		},
	}.Write(buffer)
}

func (f WAVFormat) IsFormat(data []byte) bool {

	// make sure headers match RIFF WAVE format
	_, _, _, err := RIFFHeaders(bytes.NewBuffer(data))
	return err == nil

}
