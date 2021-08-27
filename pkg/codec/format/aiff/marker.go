package aiff

import (
	"bytes"

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

func (c Marker) GetBytes() []byte {
	return utils.GetBytes(
		false,
		c.MarkerID,
		c.Position,
		c.MarkerName,
	)
}

func (c MarkerChunk) GetBytes() []byte {
	return c.MakeChunkBytes(
		c.NumMarkers,
		GetMarkersBytes(c.Markers),
	)
}

// TODO implement
func GetMarkersBytes(markers []Marker) []byte {
	return []byte("")
}

// TODO implement
func NewMarkerChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return MarkerChunk{}, nil
}
