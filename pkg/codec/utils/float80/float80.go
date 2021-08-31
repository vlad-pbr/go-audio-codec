package float80

import (
	"encoding/binary"
	"math"
	"math/big"
)

var BIAS uint16 = 16383

type Float80 struct {
	Exponent uint16
	Mantissa uint64
}

func (f Float80) UInt64() uint64 {

	// take exponent, strip sign bit, deduct bias
	exp := (f.Exponent & 32767) - BIAS

	// manually calculate mantissa float from bits
	mantissaFloat := big.NewFloat(1).SetPrec(64)
	for mask, index := uint64(1), 0; index < 64; mask, index = mask<<1, index+1 {
		if bit := int8(f.Mantissa & mask >> index); bit == 1 {
			mantissaFloat.Add(mantissaFloat, big.NewFloat(math.Pow(2, float64(-64+index))))
		}
	}

	// convert to uint64
	ans, _ := big.NewFloat(0).SetPrec(64).SetMantExp(mantissaFloat, int(exp)).Uint64()

	return ans
}

// TODO implement
func (f Float80) Bytes() []byte {
	return []byte("")
}

// TODO implement
func NewFromUInt64(value uint64) Float80 {
	return Float80{}
}

func NewFromBytes(bytes [10]byte) Float80 {

	return Float80{
		Exponent: binary.BigEndian.Uint16(bytes[0:2]),
		Mantissa: binary.BigEndian.Uint64(bytes[2:10]),
	}
}
