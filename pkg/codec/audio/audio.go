package audio

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Audio struct {
	NumChannels uint16
	SampleRate  uint64
	BitDepth    uint16
	Samples     []byte
	order       binary.ByteOrder
}

func (a Audio) String() string {
	return fmt.Sprintf("Channels: %d\nSample Rate: %d\nBit Depth: %d", a.NumChannels, a.SampleRate, a.BitDepth)
}

func (a Audio) GetSamples(order binary.ByteOrder) []byte {

	if order == a.order {
		return a.Samples
	}

	return toggleEndianness(a.BitDepth, a.Samples)
}

func NewAudio(numChannels uint16, sampleRate uint64, bitDepth uint16, samples []byte, order binary.ByteOrder) (Audio, error) {

	// init audio container
	audio := Audio{
		NumChannels: numChannels,
		SampleRate:  sampleRate,
		BitDepth:    bitDepth,
		Samples:     samples,
		order:       order,
	}

	return audio, nil

}

func toggleEndianness(bitDepth uint16, samples []byte) []byte {

	var data []byte
	sampleSize := int(math.Ceil(float64(bitDepth / 8)))

	for sampleStart, sampleEnd := 0, sampleSize-1; sampleStart < len(samples); sampleStart, sampleEnd = sampleStart+sampleSize, sampleEnd+sampleSize {
		for sample := sampleEnd; sample >= sampleStart; sample-- {
			data = append(data, samples[sample])
		}
	}

	return data
}
