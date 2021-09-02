package audio

import "fmt"

type Audio struct {
	NumChannels uint16
	NumSamples  uint64
	SampleRate  uint64
	BitDepth    uint16
	Samples     []byte
}

func (a Audio) String() string {
	return fmt.Sprintf("Channels: %d\nSamples: %d\nSample Rate: %d\nBit Depth: %d", a.NumChannels, a.NumSamples, a.SampleRate, a.BitDepth)
}
