package aiff

import (
	"bytes"
	"encoding/binary"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var NAMEID utils.FourCC = []byte("NAME")
var AUTHORID utils.FourCC = []byte("AUTH")
var COPYRIGHTID utils.FourCC = []byte("(c) ")
var ANNOTATIONID utils.FourCC = []byte("ANNO")

type TextChunk struct {
	utils.Chunk
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

func NewTextChunk(buffer *bytes.Buffer, fourCC utils.FourCC) TextChunk {

	// define chunk struct
	var textChunk TextChunk
	textChunk.ChunkID = fourCC
	textChunk.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))
	textChunk.Text = buffer.Next(int(textChunk.ChunkSize))

	return textChunk
}

func (c TextChunk) GetBytes() []byte {
	return c.GetBytesWithHeaders(c.Text)
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
