package audio

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
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

func (a *Audio) SetBitDepth(bitDepth uint16) {

	// TODO pad extra bytes if need be

	a.bitDepth = bitDepth
}

func (a *Audio) Samples(order binary.ByteOrder) []byte {

	// toggle byte order if mismatch
	if order != a.order {
		a.toggleEndianness()
	}

	return a.samples
}

func (a Audio) ByteLength() int {
	return len(a.samples)
}

func (a Audio) SamplesAmount() int {
	return a.ByteLength() / int(a.numChannels) / a.ByteDepth()
}

func (a Audio) AudioLength() float64 {
	return float64(a.SamplesAmount() / int(a.sampleRate))
}

func (a Audio) String() string {
	return fmt.Sprintf("Length: %s\n"+
		"Channels: %d\n"+
		"Sample Rate: %d\n"+
		"Bit Depth: %d", time.Duration(a.AudioLength()*float64(time.Second)), a.numChannels, a.sampleRate, a.bitDepth)
}

func (a Audio) ByteDepth() int {
	return int(math.Ceil(float64(a.bitDepth) / 8))
}

func (a *Audio) toggleEndianness() {

	byteDepth := a.ByteDepth()

	// swap order of bytes in each sample
	for sampleStart, sampleEnd := 0, byteDepth-1; sampleStart < len(a.samples); sampleStart, sampleEnd = sampleStart+byteDepth, sampleEnd+byteDepth {
		for sampleByteStart, sampleByteEnd := sampleStart, sampleEnd; sampleByteStart < sampleByteEnd; sampleByteStart, sampleByteEnd = sampleByteStart+1, sampleByteEnd-1 {
			a.samples[sampleByteStart], a.samples[sampleByteEnd] = a.samples[sampleByteEnd], a.samples[sampleByteStart]
		}
	}

	// toggle order field
	if a.order == binary.BigEndian {
		a.order = binary.LittleEndian
	} else {
		a.order = binary.BigEndian
	}
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

	// make sure samples frames fit the given samples
	{
		sampleFrameSize := audio.ByteDepth() * int(audio.NumChannels())
		samplesBytesLength := len(audio.samples)

		if samplesBytesLength%sampleFrameSize != 0 {
			return audio, fmt.Errorf("given sample bytes do not match the given sample frame size: %d (length of samples bytes array) %% %d (sample frame size) should be 0", samplesBytesLength, sampleFrameSize)
		}
	}

	return audio, nil

}
