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

// TODO implement
func NewApplicationSpecificChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return ApplicationSpecificChunk{}, nil
}
