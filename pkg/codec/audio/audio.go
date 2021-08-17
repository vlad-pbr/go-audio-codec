package audio

type Audio struct {
	NumChannels uint16
	SampleRate  uint32
	Samples     []byte
}
