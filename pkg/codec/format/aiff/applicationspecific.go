package aiff

import (
	"bytes"
	"encoding/binary"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var APPLICATIONSPECIFICID utils.FourCC = [4]byte{65, 80, 80, 76} // APPL

type ApplicationSpecificChunk struct {
	AIFFChunk
	ApplicationSignature [4]byte
	Data                 []byte
}

func (c ApplicationSpecificChunk) Write(buffer *bytes.Buffer) {
	c.WriteHeaders(buffer)
	binary.Write(buffer, binary.BigEndian, c.ApplicationSignature[:])
	binary.Write(buffer, binary.BigEndian, c.Data)
}

func NewApplicationSpecificChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	var asChunk ApplicationSpecificChunk
	asChunk.ChunkID = APPLICATIONSPECIFICID
	asChunk.ChunkSize = int32(binary.BigEndian.Uint32(utils.Next(buffer, 4)))
	copy(asChunk.ApplicationSignature[:], utils.Next(buffer, 4))
	asChunk.Data = utils.Next(buffer, int(asChunk.ChunkSize-4))

	adjustForZeroPadding(asChunk.ChunkSize, buffer)

	return ApplicationSpecificChunk{}, nil
}
