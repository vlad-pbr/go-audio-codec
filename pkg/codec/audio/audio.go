package audio

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Audio struct {
	numChannels uint16
	sampleRate  uint64
	bitDepth    uint16
	samples     []byte
	order       binary.ByteOrder
}

func (a Audio) NumChannels() uint16 {
	return a.numChannels
}

func (a Audio) SampleRate() uint64 {
	return a.sampleRate
}

func (a Audio) BitDepth() uint16 {
	return a.bitDepth
}

func (a Audio) Samples(order binary.ByteOrder) []byte {

	if order == a.order {
		return a.samples
	}

	return toggleEndianness(a.bitDepth, a.samples)
}

func (a Audio) SamplesLength() int {
	return len(a.samples)
}

func (a Audio) String() string {
	return fmt.Sprintf("Channels: %d\nSample Rate: %d\nBit Depth: %d", a.numChannels, a.sampleRate, a.bitDepth)
}

func NewAudio(numChannels uint16, sampleRate uint64, bitDepth uint16, samples []byte, order binary.ByteOrder) (Audio, error) {

	// init audio container
	audio := Audio{
		numChannels: numChannels,
		sampleRate:  sampleRate,
		bitDepth:    bitDepth,
		samples:     samples,
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
