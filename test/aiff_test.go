package test

import (
	"testing"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/format"
	"github.com/vlad-pbr/go-audio-codec/test/suite"
)

func TestIsFormat_AIFF(t *testing.T) {
	suite.IsFormat(t, format.AIFF, "samples/aiff/M1F1-int24-AFsp.aif", "samples/wav/M1F1-int24-AFsp.wav")
}

func TestDecode_AIFF(t *testing.T) {
	suite.Decode(t, format.AIFF, "samples/aiff/M1F1-int24-AFsp.aif", "samples/aiff/M1F1-int32-AFsp.aif")
}

func TestEncode_AIFF(t *testing.T) {
	suite.Encode(t, format.AIFF, format.WAV, "samples/aiff/M1F1-int24-AFsp.aif")
}
