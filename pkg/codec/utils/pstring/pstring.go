package pstring

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Pstring struct {
	length uint8
	text   []byte
}

func (p Pstring) Length() uint8 {
	return p.length
}

func (p Pstring) Text() []byte {
	return p.text
}

func (p Pstring) Bytes() []byte {

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.length)
	binary.Write(buf, binary.BigEndian, p.text)

	return buf.Bytes()
}

func New(length uint8, text []byte) (Pstring, error) {

	// make sure lengths match
	if int(length) != len(text) {
		return Pstring{}, fmt.Errorf("given pstring length does not match the given text length: length is %d while text length is %d", length, len(text))
	}

	return Pstring{
		length: length,
		text:   text,
	}, nil
}
