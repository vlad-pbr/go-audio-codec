package aiff

import (
	"bytes"
	"encoding/binary"

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

func (c Comment) Write(buffer *bytes.Buffer) {
	binary.Write(buffer, binary.BigEndian, c.TimeStamp)
	binary.Write(buffer, binary.BigEndian, c.Marker)
	binary.Write(buffer, binary.BigEndian, c.Count)
	binary.Write(buffer, binary.BigEndian, c.Text)
}

func (c CommentsChunk) Write(buffer *bytes.Buffer) {

	c.WriteHeaders(buffer)
	binary.Write(buffer, binary.BigEndian, c.NumComments)

	for _, comment := range c.Comments {
		comment.Write(buffer)
	}

}

// TODO implement
func GetCommentsBytes(comments []Comment) []byte {
	return []byte("")
}

// TODO implement
func NewCommentsChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return CommentsChunk{}, nil
}
