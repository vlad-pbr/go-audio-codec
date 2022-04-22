package suite

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/format"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec"
)

func IsFormat(t *testing.T, fi format.FormatIdentifier, correctFile string, incorrectFile string) {

	// read audio file as bytes
	data, err := ioutil.ReadFile(correctFile)
	if err != nil {
		t.Fatalf("could not read file: %s", err.Error())
	}

	// correct format
	if result := format.IdentifierToFormat[fi].IsFormat(data); !result {
		t.Errorf("got %t, expected %t", result, true)
	}

	// read audio file as bytes
	data, err = ioutil.ReadFile(incorrectFile)
	if err != nil {
		t.Fatalf("could not read file: %s", err.Error())
	}

	// incorrect format
	if result := format.IdentifierToFormat[fi].IsFormat(data); result {
		t.Errorf("wav: got %t, expected %t", result, false)
	}

}

func Decode(t *testing.T, fi format.FormatIdentifier, _24bitFile string, _32bitFile string) {

	// read audio file as bytes
	data, err := ioutil.ReadFile(_24bitFile)
	if err != nil {
		t.Fatalf("could not read file: %s", err.Error())
	}

	// decode to audio
	audio, err := codec.DecodeSpecific(data, fi)
	if err != nil {
		t.Fatalf("could not decode file: %s", err.Error())
	}

	// audio length
	if audio.Length() != 2.0 {
		t.Errorf("length: got %f, expected %f", audio.Length(), 2.0)
	}

	// number of channels
	if audio.NumChannels() != 2 {
		t.Errorf("number of channels: got %d, expected %d", audio.NumChannels(), 2)
	}

	// sample rate
	if audio.SampleRate() != 8000 {
		t.Errorf("sample rate: got %d, expected %d", audio.SampleRate(), 8000)
	}

	// bit depth
	if audio.BitDepth() != 24 {
		t.Errorf("bit depth: got %d, expected %d", audio.BitDepth(), 24)
	}

	// ========================

	// read audio file as bytes
	data, err = ioutil.ReadFile(_32bitFile)
	if err != nil {
		t.Fatalf("could not read file: %s", err.Error())
	}

	// decode to audio
	audio, err = codec.DecodeSpecific(data, fi)
	if err != nil {
		t.Fatalf("could not decode file: %s", err.Error())
	}

	// bit depth
	if audio.BitDepth() != 32 {
		t.Errorf("bit depth: got %d, expected %d", audio.BitDepth(), 32)
	}

}

func Encode(t *testing.T, from_fi format.FormatIdentifier, to_fi format.FormatIdentifier, _24bitFile string) {

	// read audio file as bytes
	data_from, err := ioutil.ReadFile(_24bitFile)
	if err != nil {
		t.Fatalf("could not read file: %s", err.Error())
	}

	// decode to audio
	audio_from, err := codec.DecodeSpecific(data_from, from_fi)
	if err != nil {
		t.Fatalf("could not decode file: %s", err.Error())
	}

	// encode to provided format
	data_to := codec.Encode(audio_from, to_fi)

	// decode to audio
	audio_to, err := codec.DecodeSpecific(data_to, to_fi)
	if err != nil {
		t.Fatalf("could not decode file: %s", err.Error())
	}

	// compare two audio containers
	if !audio_from.Equal(audio_to) {

		reason := ""

		if audio_from.String() != audio_to.String() {
			reason = fmt.Sprintf("\n\n[Original]\n%s\n\n[Reencoded]\n%s\n\n", audio_from.String(), audio_to.String())
		} else {
			reason = "playback data is correct, but audio samples differ"
		}

		t.Errorf("audio containers are not equal after reencoding: %s", reason)
	}

}
