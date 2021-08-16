package audio

type Audio struct {
	numChannels uint16
	sampleRate  uint32
	samples     []byte
}
