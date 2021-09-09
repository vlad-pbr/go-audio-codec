package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var FORMATID utils.FourCC = [4]byte{102, 109, 116, 32} // "fmt "

type FormatChunk struct {
	WAVChunk
	AudioFormat   uint16 // PCM == 1, others are currently not supported
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32 // SampleRate * NumChannels * BitsPerSample / 8
	BlockAlign    uint16 // sample size in BYTES (NumChannels * BitsPerSample / 8)
	BitsPerSample uint16 // sample size in BITS (bit depth)
}

func (c FormatChunk) Write(buffer *bytes.Buffer) {

	c.ReadHeaders(buffer)
	binary.Write(buffer, binary.LittleEndian, c.AudioFormat)
	binary.Write(buffer, binary.LittleEndian, c.NumChannels)
	binary.Write(buffer, binary.LittleEndian, c.SampleRate)
	binary.Write(buffer, binary.LittleEndian, c.ByteRate)
	binary.Write(buffer, binary.LittleEndian, c.BlockAlign)
	binary.Write(buffer, binary.LittleEndian, c.BitsPerSample)

}

func NewFormatChunk(buffer *bytes.Buffer) (FormatChunk, error) {

	// define chunk struct
	var formatChunk FormatChunk

	// set fmt chunk ID
	formatChunk.ChunkID = FORMATID
	formatChunk.ChunkSize = binary.LittleEndian.Uint32(utils.Next(buffer, 4))

	// make sure PCM is specified
	formatChunk.AudioFormat = binary.LittleEndian.Uint16(utils.Next(buffer, 2))
	if formatChunk.AudioFormat != 1 {
		return formatChunk, fmt.Errorf("only PCM audio format (1) is supported, found %d", formatChunk.AudioFormat)
	}

	formatChunk.NumChannels = binary.LittleEndian.Uint16(utils.Next(buffer, 2))
	formatChunk.SampleRate = binary.LittleEndian.Uint32(utils.Next(buffer, 4))
	formatChunk.ByteRate = binary.LittleEndian.Uint32(utils.Next(buffer, 4))
	formatChunk.BlockAlign = binary.LittleEndian.Uint16(utils.Next(buffer, 2))
	formatChunk.BitsPerSample = binary.LittleEndian.Uint16(utils.Next(buffer, 2))

	return formatChunk, nil
}
