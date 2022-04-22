package extended

import (
	"bytes"
	"encoding/binary"
	"math"
)

var BIAS uint16 = 16383

type Extended struct {
	exponent uint16 // 1st bit - sign (s), the other 15 - exponent
	mantissa uint64 // 1st bit - integer bit (i), the other 63 - mantissa
}

func NewFromBytes(bytes [10]byte) Extended {

	return Extended{
		exponent: binary.BigEndian.Uint16(bytes[0:2]),
		mantissa: binary.BigEndian.Uint64(bytes[2:10]),
	}
}

func NewFromFloat64(value float64) Extended {

	// code somewhat appropriated from
	// https://github.com/mewspring/mewmew-l/blob/c756be720bb0/internal/float80/float80.go
	// thanks, mewmew

	// get float as bits
	bits := math.Float64bits(value)

	// unpack parts
	sign := uint16(bits >> 63)
	exponent64 := uint16(bits >> 52 & 2047)
	mantissa64 := bits & 4503599627370495

	// handle zero value
	if exponent64 == 0 && mantissa64 == 0 {
		return Extended{
			exponent: sign << 15,
			mantissa: 0,
		}
	}

	// translate exponent of 64 bit float to 80 bit
	// 1023 is bias of float64
	exponent80 := exponent64 + BIAS - 1023

	// translate mantissa of 64 bit float to 80 bit
	// float64: 52 bits for mantissa
	// extended: 63 bits for mantissa
	mantissa80 := mantissa64 << (63 - 52)

	// exponent is all zeroes - zero value / subnormal
	if exponent64 == 0 {
		exponent80 = 0
	}

	// exponent is all ones - inf / NaN
	if exponent64 == 0x7FF {
		exponent80 = 0x7FFF
	}

	// align sign bit with exponent80
	exponent80 = sign<<15 | uint16(exponent80)

	// handle NaN
	if math.IsNaN(value) {
		return Extended{
			exponent: exponent80,
			mantissa: 0xC000000000000000,
		}
	}

	// align integer bit with mantissa80
	mantissa80 = 0x8000000000000000 | mantissa80

	return Extended{
		exponent: exponent80,
		mantissa: mantissa80,
	}
}

func (ext Extended) Float64() float64 {

	// unpack parts
	s := ext.exponent & uint16(32768) >> 15
	e := ext.exponent & uint16(32767)
	i := ext.mantissa & uint64(9223372036854775808) >> 63
	f := float64(0)

	// decode mantissa
	for _f, mask, exponent := ext.mantissa&uint64(9223372036854775807), uint64(4611686018427387904), 0.5; mask != 0; mask, exponent = mask>>1, exponent/2 {
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

	return math.Pow(-1, float64(s)) * math.Pow(2, float64(int64(e)-int64(BIAS))) * (float64(i) + f)
}

func (ext Extended) Bytes() [10]byte {

	buffer := new(bytes.Buffer)
	var ret [10]byte

	binary.Write(buffer, binary.BigEndian, ext.exponent)
	binary.Write(buffer, binary.BigEndian, ext.mantissa)

	copy(ret[:], buffer.Bytes())

	return ret
}
