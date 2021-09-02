package wav

import "github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"

type WAVFormat struct {
}

type riffChunk struct {
	chunkID   [4]byte // Chunk ID (should be 'RIFF')
	chunkSize uint32  // File size - 8 (without Group ID and RIFF type)
	format    [4]byte // Extension of a RIFF file (should be 'WAVE')
}

type formatChunk struct {
	chunkID       [4]byte // Chunk ID (should be 'fmt ')
	chunkSize     uint32  // Size of the rest of the chunk which follows this number
	audioFormat   uint16  // Sample format (should be 1 for 'PCM')
	numChannels   uint16  // Amount of audio channels present
	sampleRate    uint32  // Amount of samples per second of audio
	byteRate      uint32  // Amount of bytes per second of audio
	blockAlign    uint16  // Number of audio channels * Bits per Sample / 8
	bitsPerSample uint16  // Amount of bits per audio sample (bit depth)
}

type dataChunk struct {
	chunkID   [4]byte // Chunk ID (should be 'data')
	chunkSize uint32  // Number of bytes in the sample data portion
	data      []byte  // Array of audio samples
}

// TODO implement
func (f WAVFormat) Decode(bytes []byte) (audio.Audio, error) {
	return audio.Audio{}, nil
}

// TODO implement
func (f WAVFormat) Encode(audio audio.Audio) []byte {
	return []byte("")
}

// TODO implement
func (f WAVFormat) IsFormat(data []byte) bool {
	return false
}
