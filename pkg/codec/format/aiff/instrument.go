package aiff

import (
	"bytes"
	"encoding/binary"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var INSTRUMENTID utils.FourCC = [4]byte{73, 78, 83, 84} // INST

type Loop struct {
	PlayMode  int16
	BeginLoop MarkerID
	EndLoop   MarkerID
}

type InstrumentChunk struct { // size is always 20
	AIFFChunk
	BaseNote     int8
	Detune       int8
	LowNote      int8
	HighNote     int8
	LowVelocity  int8
	HighVelocity int8
	Gain         int16
	SustainLoop  Loop
	ReleaseLoop  Loop
}

func (l Loop) Write(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, l.PlayMode)
	binary.Write(buffer, binary.BigEndian, l.BeginLoop)
	binary.Write(buffer, binary.BigEndian, l.EndLoop)
}

func (c InstrumentChunk) Write(buffer *bytes.Buffer) {
	c.WriteHeaders(buffer)
	binary.Write(buffer, binary.BigEndian, c.BaseNote)
	binary.Write(buffer, binary.BigEndian, c.Detune)
	binary.Write(buffer, binary.BigEndian, c.LowNote)
	binary.Write(buffer, binary.BigEndian, c.HighNote)
	binary.Write(buffer, binary.BigEndian, c.LowVelocity)
	binary.Write(buffer, binary.BigEndian, c.HighVelocity)
	binary.Write(buffer, binary.BigEndian, c.Gain)
	c.SustainLoop.Write(buffer)
	c.ReleaseLoop.Write(buffer)
}

// TODO implement
func NewInstrumentChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return InstrumentChunk{}, nil
}
