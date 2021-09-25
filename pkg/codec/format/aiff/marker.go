package aiff

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils/pstring"
)

var MARKERID utils.FourCC = [4]byte{77, 65, 82, 75} // MARK

type MarkerID int16

type Marker struct {
	MarkerID   MarkerID
	Position   uint32
	MarkerName pstring.Pstring
}

type MarkerChunk struct {
	AIFFChunk
	NumMarkers uint16
	Markers    []Marker
}

func (c Marker) Write(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, c.MarkerID)
	binary.Write(buffer, binary.BigEndian, c.Position)
	buffer.Write(c.MarkerName.Bytes())

	adjustForZeroPadding(int32(c.MarkerName.Length()), buffer, true)
}

func (c MarkerChunk) Write(buffer *bytes.Buffer) {

	c.WriteHeaders(buffer)
	binary.Write(buffer, binary.BigEndian, c.NumMarkers)

	for _, marker := range c.Markers {
		marker.Write(buffer)
	}

}

func NewMarkerChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	var markerChunk MarkerChunk
	markerChunk.ChunkID = MARKERID
	markerChunk.ChunkSize = int32(binary.BigEndian.Uint32(utils.Next(buffer, 4)))
	markerChunk.NumMarkers = binary.BigEndian.Uint16(utils.Next(buffer, 2))

	for i := 0; i < int(markerChunk.NumMarkers); i++ {

		// init marker
		marker := Marker{
			MarkerID: MarkerID(binary.BigEndian.Uint16(utils.Next(buffer, 2))),
			Position: binary.BigEndian.Uint32(utils.Next(buffer, 4)),
		}

		// decode pstring
		pstringLength := uint8(utils.Next(buffer, 1)[0])
		var err error
		marker.MarkerName, err = pstring.New(pstringLength, utils.Next(buffer, int(pstringLength)))

		if err != nil {
			return markerChunk, fmt.Errorf("could not decode marker name: %s", err.Error())
		}

		// adjust by pstring length
		adjustForZeroPadding(int32(pstringLength), buffer, false)
	}

	return markerChunk, nil
}
