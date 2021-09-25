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

func NewCommentsChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {

	// define chunk struct
	var commentsChunk CommentsChunk
	commentsChunk.ChunkID = COMMENTSID
	commentsChunk.ChunkSize = int32(binary.BigEndian.Uint32(utils.Next(buffer, 4)))
	commentsChunk.NumComments = binary.BigEndian.Uint16(utils.Next(buffer, 2))

	for i := 0; i < int(commentsChunk.NumComments); i++ {

		// read comment chunk
		comment := Comment{
			TimeStamp: binary.BigEndian.Uint32(utils.Next(buffer, 4)),
			Marker:    MarkerID(binary.BigEndian.Uint16(utils.Next(buffer, 2))),
			Count:     binary.BigEndian.Uint16(utils.Next(buffer, 2)),
		}
		comment.Text = utils.Next(buffer, int(comment.Count))

		// adjust by count
		adjustForZeroPadding(int32(comment.Count), buffer)

		// add to comments slice
		commentsChunk.Comments = append(commentsChunk.Comments, comment)
	}

	return commentsChunk, nil
}
