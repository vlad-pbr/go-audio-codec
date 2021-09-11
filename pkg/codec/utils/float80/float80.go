package float80

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/big"
)

var BIAS uint16 = 16383

type Float80 struct {
	Exponent uint16
	Mantissa uint64
}

func (f Float80) Float() *big.Float {

	// take exponent, strip sign bit, deduct bias
	exp := (f.Exponent & 32767) - BIAS

	// manually calculate mantissa float from bits
	mantissaFloat := big.NewFloat(1).SetPrec(64)
	for mask, index := uint64(1), 0; index < 64; mask, index = mask<<1, index+1 {
		if bit := int8(f.Mantissa & mask >> index); bit == 1 {
			mantissaFloat.Add(mantissaFloat, big.NewFloat(math.Pow(2, float64(-64+index))))
		}
	}

	return big.NewFloat(0).SetPrec(64).SetMantExp(mantissaFloat, int(exp))
}

func (f Float80) Bytes() [10]byte {

	buffer := new(bytes.Buffer)
	var ret [10]byte

	binary.Write(buffer, binary.BigEndian, f.Exponent)
	binary.Write(buffer, binary.BigEndian, f.Mantissa)

	copy(ret[:], buffer.Bytes())

	return ret
}

// TODO implement
func NewFromFloat(value *big.Float) Float80 {
	return Float80{}
}

func NewFromBytes(bytes [10]byte) Float80 {

	return Float80{
		Exponent: binary.BigEndian.Uint16(bytes[0:2]),
		Mantissa: binary.BigEndian.Uint64(bytes[2:10]),
	}
}
