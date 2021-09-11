package aiff

import (
	"bytes"
	"encoding/binary"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var MARKERID utils.FourCC = [4]byte{77, 65, 82, 75} // MARK

// TODO pascal string
type pstring []byte

type MarkerID int16

type Marker struct {
	MarkerID   MarkerID
	Position   uint32
	MarkerName pstring
}

type MarkerChunk struct {
	AIFFChunk
	NumMarkers uint16
	Markers    []Marker
}

func (c Marker) Write(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, c.MarkerID)
	binary.Write(buffer, binary.BigEndian, c.Position)
	binary.Write(buffer, binary.BigEndian, c.MarkerName)
}

func (c MarkerChunk) Write(buffer *bytes.Buffer) {

	c.WriteHeaders(buffer)
	binary.Write(buffer, binary.BigEndian, c.NumMarkers)

	for _, marker := range c.Markers {
		marker.Write(buffer)
	}

}

// TODO implement
func NewMarkerChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	var markerChunk MarkerChunk
	markerChunk.ChunkID = MARKERID
	markerChunk.ChunkSize = int32(binary.BigEndian.Uint32(utils.Next(buffer, 4)))
	markerChunk.NumMarkers = binary.BigEndian.Uint16(utils.Next(buffer, 2))

	// TODO

	return MarkerChunk{}, nil
}
