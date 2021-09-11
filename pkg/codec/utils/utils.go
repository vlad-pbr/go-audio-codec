package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type FourCC [4]byte

func (f FourCC) GetBytes() []byte {
	return f[:]
}

type Chunk struct {
	ChunkID FourCC
}

func (c Chunk) GetBytesWithID(fields ...interface{}) []byte {
	return append(GetBytes(false, c.ChunkID), GetBytes(false, fields...)...)
}

type ChunkInterface interface {
	GetID() FourCC
	Write(buffer *bytes.Buffer)
}

func (c Chunk) GetID() FourCC {
	return c.ChunkID
}

func GetBytes(zeroPad bool, fields ...interface{}) []byte {

	var output []byte
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	for _, field := range fields {

		// convert interface to bytes
		if err := encoder.Encode(field); err != nil {
			panic(err)
		}

		output = append(output, buffer.Bytes()...)
	}

	// zero pad odd output if specified
	if zeroPad && len(output)%2 != 0 {
		output = append(output, byte(0))
	}

	return output
}

func ContainsFourCC(slice []FourCC, fourCC FourCC) bool {
	for _, item := range slice {
		if bytes.Equal(item[:], fourCC[:]) {
			return true
		}
	}

	return false
}

// TODO implement
func GetChunksBytes(chunks []ChunkInterface) []byte {
	return []byte("")
}

// panicing version of buffer.Next
func Next(buffer *bytes.Buffer, amount int) []byte {

	if buffer.Len() < amount {
		panic(fmt.Sprintf("unexpected EOF: tried to read %d bytes when only %d left in buffer", amount, buffer.Len()))
	}

	return buffer.Next(amount)
}
