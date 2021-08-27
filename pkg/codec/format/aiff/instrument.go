package aiff

import (
	"bytes"

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

func (l Loop) GetBytes() []byte {
	return utils.GetBytes(
		false,
		l.PlayMode,
		l.BeginLoop,
		l.EndLoop,
	)
}

func (c InstrumentChunk) GetBytes() []byte {
	return c.MakeChunkBytes(
		c.BaseNote,
		c.Detune,
		c.LowNote,
		c.HighNote,
		c.LowVelocity,
		c.HighVelocity,
		c.Gain,
		c.SustainLoop.GetBytes(),
		c.ReleaseLoop.GetBytes(),
	)
}

// TODO implement
func NewInstrumentChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return InstrumentChunk{}, nil
}
