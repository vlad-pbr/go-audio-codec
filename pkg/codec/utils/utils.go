package utils

import (
	"bytes"
	"encoding/gob"
)

type FourCC []byte

type Chunk struct {
	ChunkID FourCC
	// ChunkSize int32
	// for AIFF, data must be zero padded if odd length
}

func (c Chunk) GetBytesWithID(fields ...interface{}) []byte {
	return append(GetBytes(c.ChunkID), GetBytes(fields...)...)
}

type ChunkInterface interface {
	GetID() FourCC
	// GetSize() int32
	GetBytes() []byte
}

func (c Chunk) GetID() FourCC {
	return c.ChunkID
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

func ContainsFourCC(slice []FourCC, fourCC FourCC) bool {
	for _, item := range slice {
		if bytes.Compare(item, fourCC) == 0 {
			return true
		}
	}

	return false
}

// TODO implement
func GetChunksBytes(chunks []ChunkInterface) []byte {
	return []byte("")
}
