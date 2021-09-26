package float80

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

var BIAS uint16 = 16383

type Float80 struct {
	Exponent uint16
	Mantissa uint64
}

func (f Float80) Float() *big.Float {

	float := big.NewFloat(0).SetPrec(64)

	// iterate bits, then sum accordingly to big float
	for mask, index, exponent := uint64(1), 0, int(f.Exponent-BIAS)-63; index < 64; mask, index, exponent = mask<<1, index+1, exponent+1 {
		if bit := int8(f.Mantissa & mask >> index); bit == 1 {
			if exponent >= 0 {
				float.Add(float, big.NewFloat(math.Pow(2, float64(exponent))))
			} else {
				float.Add(float, big.NewFloat(1.0/math.Pow(2, float64(exponent))))
			}
		}
	}

	return float
}

// TODO implement
func (f Float80) Float64() float64 {
	panic("not implemented")
}

func (f Float80) Bytes() [10]byte {

	buffer := new(bytes.Buffer)
	var ret [10]byte

	binary.Write(buffer, binary.BigEndian, f.Exponent)
	binary.Write(buffer, binary.BigEndian, f.Mantissa)

	copy(ret[:], buffer.Bytes())

	return ret
}

func NewFromFloat(value *big.Float) Float80 {

	// big.Float annoyingly does not have a proper getter for the mantissa in any format
	// so we format it through text and parse the decimal -_-

	// NOTE: this is temporary as only positive values are currently supported

	// ensure value has correct precision
	preciseValue := new(big.Float).SetPrec(64).Set(value)

	// get decimal representation of mantissa in text
	text := preciseValue.Text(byte(98), 0)

	// find index of 'p' char located after mantissa
	var pindex int
	for pindex = 0; text[pindex] != byte(112) && pindex < len(text); pindex++ {
	}

	// parse mantissa to uint64
	mantissa, err := strconv.ParseUint(text[:pindex], 10, 64)
	if err != nil {
		panic(fmt.Errorf("could not parse decimal mantissa: %s", err.Error()))
	}

	return Float80{
		Exponent: uint16(preciseValue.MantExp(nil)) + BIAS - 1,
		Mantissa: mantissa,
	}
}

// TODO implement
func NewFromFloat64(value float64) Float80 {
	panic("not implemented")
}

func NewFromBytes(bytes [10]byte) Float80 {

	return Float80{
		Exponent: binary.BigEndian.Uint16(bytes[0:2]),
		Mantissa: binary.BigEndian.Uint64(bytes[2:10]),
	}
}
