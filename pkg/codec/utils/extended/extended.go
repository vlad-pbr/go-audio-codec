package extended

import (
	"bytes"
	"encoding/binary"
	"math"
)

var BIAS uint16 = 16383

type Extended struct {
	left  uint16 // 1st bit - sign (s), the other 15 - exponent
	right uint64 // 1st bit - i, the other 63 - mantissa
}

func NewFromBytes(bytes [10]byte) Extended {

	return Extended{
		left:  binary.BigEndian.Uint16(bytes[0:2]),
		right: binary.BigEndian.Uint64(bytes[2:10]),
	}
}

// TODO
func NewFromFloat64(value float64) Extended {

	return Extended{
		left:  0,
		right: 0,
	}
}

func (ext Extended) Float64() float64 {

	// unpack parts
	s := ext.left & uint16(32768) >> 15
	e := ext.left & uint16(32767)
	i := ext.right & uint64(9223372036854775808) >> 63
	f := float64(0)

	// decode mantissa
	for _f, mask, exponent := ext.right&uint64(9223372036854775807), uint64(4611686018427387904), 0.5; mask != 0; mask, exponent = mask>>1, exponent/2 {
		if _f&mask > 0 {
			f += exponent
		}
	}

	// handle +-Inf
	if e == 32767 && f == 0 {
		return math.Pow(-1, float64(s)) * math.Inf(1)
	}

	// handle NaN
	if e == 32767 && f != 0 {
		return math.NaN()
	}

	return math.Pow(-1, float64(s)) * math.Pow(2, float64(e-BIAS)) * (float64(i) + f)
}

func (ext Extended) Bytes() [10]byte {

	buffer := new(bytes.Buffer)
	var ret [10]byte

	binary.Write(buffer, binary.BigEndian, ext.left)
	binary.Write(buffer, binary.BigEndian, ext.right)

	copy(ret[:], buffer.Bytes())

	return ret
}
