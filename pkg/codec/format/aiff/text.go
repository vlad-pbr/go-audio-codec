package aiff

import (
	"bytes"
	"encoding/binary"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var NAMEID utils.FourCC = [4]byte{78, 65, 77, 69}       // NAME
var AUTHORID utils.FourCC = [4]byte{65, 85, 84, 72}     // AUTH
var COPYRIGHTID utils.FourCC = [4]byte{40, 99, 41, 32}  // (c)
var ANNOTATIONID utils.FourCC = [4]byte{65, 78, 78, 79} // ANNO

type TextChunk struct {
	AIFFChunk
	Text []byte
}

type NameChunk struct {
	TextChunk
}

type AuthorChunk struct {
	TextChunk
}

type CopyrightChunk struct {
	TextChunk
}

type AnnotationChunk struct {
	TextChunk
}

func (c TextChunk) GetBytes() []byte {
	return c.MakeChunkBytes(
		c.Text,
	)
}

func NewTextChunk(buffer *bytes.Buffer, fourCC utils.FourCC) TextChunk {

	// define chunk struct
	var textChunk TextChunk
	textChunk.ChunkID = fourCC
	textChunk.ChunkSize = int32(binary.BigEndian.Uint32(utils.Next(buffer, 4)))
	textChunk.Text = utils.Next(buffer, int(textChunk.ChunkSize))

	AdjustForZeroPadding(textChunk.ChunkSize, buffer)

	return textChunk
}

func NewNameChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return NameChunk{NewTextChunk(buffer, NAMEID)}, nil
}

func NewAuthorChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return AuthorChunk{NewTextChunk(buffer, AUTHORID)}, nil
}

func NewCopyrightChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return CopyrightChunk{NewTextChunk(buffer, COPYRIGHTID)}, nil
}

func NewAnnotationChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return AnnotationChunk{NewTextChunk(buffer, ANNOTATIONID)}, nil
}
