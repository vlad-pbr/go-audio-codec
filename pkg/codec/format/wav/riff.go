package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

type RIFFChunk struct {
	WAVChunk
	Format      utils.FourCC
	FormatChunk FormatChunk
	DataChunk   DataChunk
}

var RIFFID utils.FourCC = [4]byte{82, 73, 70, 70} // RIFF
var WAVEID utils.FourCC = [4]byte{87, 65, 86, 69} // WAVE

func (c RIFFChunk) GetBytes() []byte {
	return c.MakeChunkBytes(
		c.Format,
	)
}

func RIFFHeaders(buffer *bytes.Buffer) (utils.FourCC, uint32, utils.FourCC, error) {

	var chunkID utils.FourCC
	var chunkSize uint32
	var format utils.FourCC

	// parse riff chunk ID
	copy(chunkID[:], utils.Next(buffer, 4))
	if !bytes.Equal(chunkID[:], RIFFID[:]) {
		return chunkID, chunkSize, format, fmt.Errorf("RIFF chunk ID is invalid: found %s, must be %s", chunkID, RIFFID)
	}

	// parse riff chunk size
	chunkSize = binary.LittleEndian.Uint32(utils.Next(buffer, 4))

	// parse riff type
	copy(format[:], utils.Next(buffer, 4))
	if !bytes.Equal(format[:], WAVEID[:]) {
		return chunkID, chunkSize, format, fmt.Errorf("RIFF format is invalid: found %s, must be %s", format, WAVEID)
	}

	return chunkID, chunkSize, format, nil
}

func NewRIFFChunk(buffer *bytes.Buffer) (RIFFChunk, error) {

	var riff RIFFChunk
	var err error

	// parse riff chunk headers
	riff.ChunkID, riff.ChunkSize, riff.Format, err = RIFFHeaders(buffer)
	if err != nil {
		return riff, fmt.Errorf("error while decoding RIFF chunk headers: %s", err.Error())
	}

	// decode fmt
	riff.FormatChunk, err = NewFormatChunk(buffer)
	if err != nil {
		return riff, fmt.Errorf("error while decoding fmt chunk: %s", err.Error())
	}

	// decode data
	riff.DataChunk, err = NewDataChunk(buffer)
	if err != nil {
		return riff, fmt.Errorf("error while decoding data chunk: %s", err.Error())
	}

	return riff, nil
}
