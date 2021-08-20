package aiff

import "github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"

var MARKERID utils.FourCC = []byte("MARK")

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
	return utils.GetBytes(c.MarkerID, c.Position, c.MarkerName)
}

func (c MarkerChunk) GetBytes() []byte {
	return c.GetBytesWithHeaders(c.NumMarkers, GetMarkersBytes(c.Markers))
}

// TODO implement
func GetMarkersBytes(markers []Marker) []byte {
	return []byte("")
}
