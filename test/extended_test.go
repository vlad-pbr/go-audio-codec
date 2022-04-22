package test

import (
	"math"
	"testing"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils/extended"
)

func TestBytesToFloat64(t *testing.T) {
	if result := extended.NewFromBytes([10]byte{64, 11, 250, 0, 0, 0, 0, 0, 0, 0}).Float64(); result != 8000.0 {
		t.Errorf("got %f, expected %f", result, 8000.0)
	}
}

func TestFloat64ToBytes(t *testing.T) {
	if result := extended.NewFromFloat64(8000.0).Bytes(); result != [10]byte{64, 11, 250, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("got %b, expected %b", result, [10]byte{64, 11, 250, 0, 0, 0, 0, 0, 0, 0})
	}
}

func TestBytesToFloat64ToBytes(t *testing.T) {
	if result := extended.NewFromFloat64(extended.NewFromBytes([10]byte{64, 11, 250, 0, 0, 0, 0, 0, 0, 0}).Float64()).Bytes(); result != [10]byte{64, 11, 250, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("got %b, expected %b", result, [10]byte{64, 11, 250, 0, 0, 0, 0, 0, 0, 0})
	}
}

func TestFloat64ToBytesToFloat64(t *testing.T) {
	if result := extended.NewFromBytes(extended.NewFromFloat64(8000.0).Bytes()).Float64(); result != 8000.0 {
		t.Errorf("got %f, expected %f", result, 8000.0)
	}
}

func TestBytesValues(t *testing.T) {

	// -8000
	if result := extended.NewFromBytes([10]byte{192, 11, 250, 0, 0, 0, 0, 0, 0, 0}).Float64(); result != -8000.0 {
		t.Errorf("got %f, expected %f", result, -8000.0)
	}

	// -8000.125
	if result := extended.NewFromBytes([10]byte{192, 11, 250, 1, 0, 0, 0, 0, 0, 0}).Float64(); result != -8000.125 {
		t.Errorf("got %f, expected %f", result, -8000.125)
	}

	// -144117439875669557248.0
	if result := extended.NewFromBytes([10]byte{192, 65, 250, 1, 0, 0, 0, 0, 0, 0}).Float64(); result != -144117439875669557248.0 {
		t.Errorf("got %f, expected %f", result, -144117439875669557248.0)
	}

	// +Inf
	if result := extended.NewFromBytes([10]byte{127, 255, 128, 0, 0, 0, 0, 0, 0, 0}).Float64(); result != math.Inf(1) {
		t.Errorf("got %f, expected %f", result, math.Inf(1))
	}

	// -Inf
	if result := extended.NewFromBytes([10]byte{255, 255, 128, 0, 0, 0, 0, 0, 0, 0}).Float64(); result != math.Inf(-1) {
		t.Errorf("got %f, expected %f", result, math.Inf(-1))
	}

	// NaN
	if result := extended.NewFromBytes([10]byte{127, 255, 192, 0, 0, 0, 0, 0, 0, 0}).Float64(); !math.IsNaN(result) {
		t.Errorf("got %f, expected %f", result, math.NaN())
	}

}

func TestFloat64Values(t *testing.T) {

	// -8000
	if result := extended.NewFromFloat64(-8000.0).Bytes(); result != [10]byte{192, 11, 250, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("got %b, expected %b", result, [10]byte{192, 11, 250, 0, 0, 0, 0, 0, 0, 0})
	}

	// -8000.125
	if result := extended.NewFromFloat64(-8000.125).Bytes(); result != [10]byte{192, 11, 250, 1, 0, 0, 0, 0, 0, 0} {
		t.Errorf("got %b, expected %b", result, [10]byte{192, 11, 250, 1, 0, 0, 0, 0, 0, 0})
	}

	// -144117439875669557248.0
	if result := extended.NewFromFloat64(-144117439875669557248.0).Bytes(); result != [10]byte{192, 65, 250, 1, 0, 0, 0, 0, 0, 0} {
		t.Errorf("got %b, expected %b", result, [10]byte{192, 65, 250, 1, 0, 0, 0, 0, 0, 0})
	}

	// +Inf
	if result := extended.NewFromFloat64(math.Inf(1)).Bytes(); result != [10]byte{127, 255, 128, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("got %b, expected %b", result, [10]byte{127, 255, 128, 0, 0, 0, 0, 0, 0, 0})
	}

	// -Inf
	if result := extended.NewFromFloat64(math.Inf(-1)).Bytes(); result != [10]byte{255, 255, 128, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("got %b, expected %b", result, [10]byte{127, 255, 128, 0, 0, 0, 0, 0, 0, 0})
	}

	// NaN
	if result := extended.NewFromFloat64(math.NaN()).Bytes(); result != [10]byte{127, 255, 192, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("got %b, expected %b", result, [10]byte{127, 255, 128, 0, 0, 0, 0, 0, 0, 0})
	}
}
