package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var DATAID utils.FourCC = [4]byte{100, 97, 116, 97} // data

type DataChunk struct {
	WAVChunk
	Data []byte
}

func (c DataChunk) GetBytes() []byte {
	return c.MakeChunkBytes()
}

func NewDataChunk(buffer *bytes.Buffer) (DataChunk, error) {

	// define chunk struct
	var dataChunk DataChunk

	// parse data chunk ID
	copy(dataChunk.ChunkID[:], utils.Next(buffer, 4))
	if !bytes.Equal(dataChunk.ChunkID[:], DATAID[:]) {
		return dataChunk, fmt.Errorf("data chunk ID is invalid: found %s, must be %s", dataChunk.ChunkID, DATAID)
	}

	dataChunk.ChunkSize = binary.LittleEndian.Uint32(utils.Next(buffer, 4))
	dataChunk.Data = utils.Next(buffer, int(dataChunk.ChunkSize)) // TODO this is terrible

	return dataChunk, nil
}
