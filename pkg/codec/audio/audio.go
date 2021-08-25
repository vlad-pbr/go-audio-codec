package audio

type Audio struct {
	NumChannels uint16
	SampleRate  uint32
	BitDepth    uint16
	BitRate     uint32
	Samples     []byte
}
