package utils

type FourCC []byte

type Chunk struct {
	ChunkID   FourCC
	ChunkSize int32
}

type ChunkInterface interface {
	GetID() FourCC
	GetSize() int32
	GetData() []byte
	GetBytes() []byte
}

func (c Chunk) GetID() FourCC {
	return c.ChunkID
}

func (c Chunk) GetSize() int32 {
	return c.ChunkSize
}

// TODO generic GetBytes for chunk

func GetBytes(fields ...interface{}) []byte {
	// TODO convert all interfaces to single byte array, return
	return []byte("placeholder")
}
