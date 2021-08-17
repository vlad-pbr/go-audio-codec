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

// TODO return proper chunk instead of generic TextChunk

func NewTextChunk(buffer *bytes.Buffer, fourCC utils.FourCC) (utils.ChunkInterface, error) {

	// define chunk struct
	var textChunk TextChunk
	textChunk.ChunkID = fourCC
	textChunk.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))
	textChunk.Text = buffer.Next(int(textChunk.ChunkSize))

	return textChunk, nil
}

func (c TextChunk) GetBytes() []byte {
	return c.GetBytesWithHeaders(c.Text)
}

type NameChunk struct {
	TextChunk
}

func NewNameChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return NewTextChunk(buffer, NAMEID)
}

type AuthorChunk struct {
	TextChunk
}

func NewAuthorChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return NewTextChunk(buffer, AUTHORID)
}

type CopyrightChunk struct {
	TextChunk
}

func NewCopyrightChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return NewTextChunk(buffer, COPYRIGHTID)
}

type AnnotationChunk struct {
	TextChunk
}

func NewAnnotationChunk(buffer *bytes.Buffer) (utils.ChunkInterface, error) {
	return NewTextChunk(buffer, ANNOTATIONID)
}
