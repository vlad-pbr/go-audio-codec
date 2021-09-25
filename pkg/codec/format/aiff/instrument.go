package aiff

import (
	"bytes"
	"encoding/binary"
	"fmt"

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

func NewInstrumentChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	// define chunk struct
	var instChunk InstrumentChunk
	instChunk.ChunkID = INSTRUMENTID
	instChunk.ChunkSize = int32(binary.BigEndian.Uint32(utils.Next(buffer, 4)))

	// make sure chunk size is 20
	if instChunk.ChunkSize != 18 {
		return instChunk, fmt.Errorf("%s chunk size is invalid: found %d, must be %d", string(INSTRUMENTID[:]), instChunk.ChunkSize, 20)
	}

	instChunk.BaseNote = int8(utils.Next(buffer, 1)[0])
	instChunk.Detune = int8(utils.Next(buffer, 1)[0])
	instChunk.LowNote = int8(utils.Next(buffer, 1)[0])
	instChunk.HighNote = int8(utils.Next(buffer, 1)[0])
	instChunk.LowVelocity = int8(utils.Next(buffer, 1)[0])
	instChunk.HighVelocity = int8(utils.Next(buffer, 1)[0])
	instChunk.Gain = int16(binary.BigEndian.Uint16(utils.Next(buffer, 2)))

	instChunk.SustainLoop = Loop{
		PlayMode:  int16(binary.BigEndian.Uint16(utils.Next(buffer, 2))),
		BeginLoop: MarkerID(binary.BigEndian.Uint16(utils.Next(buffer, 2))),
		EndLoop:   MarkerID(binary.BigEndian.Uint16(utils.Next(buffer, 2))),
	}

	instChunk.ReleaseLoop = Loop{
		PlayMode:  int16(binary.BigEndian.Uint16(utils.Next(buffer, 2))),
		BeginLoop: MarkerID(binary.BigEndian.Uint16(utils.Next(buffer, 2))),
		EndLoop:   MarkerID(binary.BigEndian.Uint16(utils.Next(buffer, 2))),
	}

	return instChunk, nil
}
