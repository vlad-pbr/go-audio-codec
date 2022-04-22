package test

import (
	"testing"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/format"
	"github.com/vlad-pbr/go-audio-codec/test/suite"
)

func TestIsFormat_WAV(t *testing.T) {
	suite.IsFormat(t, format.WAV, "samples/wav/M1F1-int24-AFsp.wav", "samples/aiff/M1F1-int24-AFsp.aif")
}

func TestDecode_WAV(t *testing.T) {
	suite.Decode(t, format.WAV, "samples/wav/M1F1-int24-AFsp.wav", "samples/wav/M1F1-int32-AFsp.wav")
}

func TestEncode_WAV(t *testing.T) {
	suite.Encode(t, format.WAV, format.AIFF, "samples/wav/M1F1-int24-AFsp.wav")
}
