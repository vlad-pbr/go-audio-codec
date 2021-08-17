package utils

import (
	"bytes"
	"encoding/gob"
)

type FourCC []byte

type Chunk struct {
	ChunkID   FourCC
	ChunkSize int32
}

func (c Chunk) GetBytesWithHeaders(fields ...interface{}) []byte {
	return append(GetBytes(c.ChunkID, c.ChunkSize), GetBytes(fields...)...)
}

type ChunkInterface interface {
	GetID() FourCC
	GetSize() int32
	GetBytes() []byte
}

func (c Chunk) GetID() FourCC {
	return c.ChunkID
}

func (c Chunk) GetSize() int32 {
	return c.ChunkSize
}

func GetBytes(fields ...interface{}) []byte {

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

	return output
}
