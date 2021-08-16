package aiff

import (
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var FORMID utils.FourCC = []byte("FORM")
var COMMONID utils.FourCC = []byte("COMM")
var SOUNDID utils.FourCC = []byte("SSND")

// TODO float80
type extended [10]byte

type AIFFFormat struct {
}

type formChunk struct {
	utils.Chunk
	formType    utils.FourCC
	localChunks []utils.ChunkInterface
}

// TODO multichannel table

type commonChunk struct { // size is always 18, one and ONLY one required
	utils.Chunk
	numChannels     int16
	numSampleFrames uint32
	sampleSize      int16
	sampleRate      extended
}

func (c commonChunk) GetData() []byte {
	return utils.GetBytes(c.numChannels, c.numSampleFrames, c.sampleSize, c.sampleRate)
}

type soundDataChunk struct { // one and only one UNLESS numSampleFrames is 0, then not required
	utils.Chunk
	offset    uint32
	blockSize uint32
	soundData []uint8
}

func (c soundDataChunk) GetData() []byte {
	return utils.GetBytes(c.offset, c.blockSize, c.soundData)
}

// TODO optional chunks
