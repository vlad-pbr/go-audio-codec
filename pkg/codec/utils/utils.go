package utils

import (
	"bytes"
	"fmt"
)

type FourCC [4]byte

func (f FourCC) GetBytes() []byte {
	return f[:]
}

type Chunk struct {
	ChunkID FourCC
}

type ChunkInterface interface {
	GetID() FourCC
	Write(buffer *bytes.Buffer)
}

func (c Chunk) GetID() FourCC {
	return c.ChunkID
}

func ContainsFourCC(slice []FourCC, fourCC FourCC) bool {
	for _, item := range slice {
		if bytes.Equal(item[:], fourCC[:]) {
			return true
		}
	}

	return false
}

// panicing version of buffer.Next
func Next(buffer *bytes.Buffer, amount int) []byte {

	if buffer.Len() < amount {
		panic(fmt.Sprintf("unexpected EOF: tried to read %d bytes when only %d left in buffer", amount, buffer.Len()))
	}

	return buffer.Next(amount)
}
