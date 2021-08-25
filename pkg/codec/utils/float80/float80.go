package float80

import (
	"encoding/binary"
)

type Float80 struct {
	Exponent uint16
	Mantissa uint64
}

func NewFromBytes(bytes [10]byte) Float80 {

	return Float80{
		Exponent: binary.BigEndian.Uint16(bytes[0:2]),
		Mantissa: binary.BigEndian.Uint64(bytes[2:10]),
	}
}

// TODO implement
func NewFromUInt64(integer uint64) Float80 {
	return Float80{}
}

// TODO implement
func (f Float80) Multiply(multiplier Float80) Float80 {
	return multiplier
}

// TODO implement
func (f Float80) ToBytes() []byte {
	return []byte("")
}
