package aiff

import (
	"bytes"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var COMMENTSID utils.FourCC = [4]byte{67, 79, 77, 84} // COMT

type Comment struct {
	TimeStamp uint32
	Marker    MarkerID
	Count     uint16
	Text      []byte
}

type CommentsChunk struct {
	AIFFChunk
	NumComments uint16
	Comments    []Comment
}

func (c Comment) GetBytes() []byte {
	return utils.GetBytes(
		true,
		c.TimeStamp,
		c.Marker,
		c.Count,
		c.Text,
	)
}

func (c CommentsChunk) GetBytes() []byte {
	return c.MakeChunkBytes(
		c.NumComments,
		GetCommentsBytes(c.Comments),
	)
}

// TODO implement
func GetCommentsBytes(comments []Comment) []byte {
	return []byte("")
}

// TODO implement
func NewCommentsChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return CommentsChunk{}, nil
}
