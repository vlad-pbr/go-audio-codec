package aiff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var FORMID utils.FourCC = []byte("FORM")
var FORMTYPE utils.FourCC = []byte("AIFF")

// possible local chunks, in order of precedence
var LocalChunks = map[string]interface{}{
	string(COMMONID): NewCommonChunk,
	string(SOUNDID):  NewSoundChunk,
}

// local chunks which are required in a valid AIFF file
var RequiredLocalChunkIDs = []utils.FourCC{
	COMMONID,
	SOUNDID,
}

// chunks which are allowed to be present more than once
var AllowedMultipleChunks = []utils.FourCC{}

type FormChunk struct {
	utils.Chunk
	formType    utils.FourCC // must be AIFF
	localChunks []utils.ChunkInterface
}

func NewFormChunk(buffer *bytes.Buffer) (FormChunk, error) {

	var form FormChunk

	// parse form chunk ID
	form.ChunkID = buffer.Next(4)
	if bytes.Compare(form.ChunkID, FORMID) != 0 {
		return FormChunk{}, errors.New(fmt.Sprintf("FORM chunk ID is invalid: found %s, must be %s", form.ChunkID, FORMID))
	}

	// parse form chunk size
	form.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))

	// parse form chunk ID
	form.formType = buffer.Next(4)
	if bytes.Compare(form.formType, FORMTYPE) != 0 {
		return FormChunk{}, errors.New(fmt.Sprintf("FORM chunk type is invalid: found %s, must be %s", form.formType, FORMTYPE))
	}

	// the following chunks can be present
	var presentChunks = map[string]bool{}
	for fourCC := range LocalChunks {
		presentChunks[fourCC] = false
	}

	// while unread portion of buffer is not empty
	for buffer.Len() != 0 {

		chunkID := buffer.Next(4)

		// make sure chunk not already present unless allowed to be present multiple times
		if presentChunks[string(chunkID)] && !utils.ContainsFourCC(AllowedMultipleChunks, chunkID) {
			return FormChunk{}, errors.New(fmt.Sprintf("more than one instance of %s chunk present", string(chunkID)))
		}

		// retrieve target chunk's creation function
		newChunkFunction := LocalChunks[string(chunkID)]
		if newChunkFunction == nil {
			return FormChunk{}, errors.New(fmt.Sprintf("invalid chunk ID found: %s", string(chunkID)))
		}

		// create target chunk
		chunk, err := newChunkFunction.(func(*bytes.Buffer) (utils.ChunkInterface, error))(buffer)
		if err != nil {
			return FormChunk{}, errors.New(fmt.Sprintf("error occurred while decoding %s chunk: %s", string(chunkID), err.Error()))
		}

		// append resulting local chunk
		form.localChunks = append(form.localChunks, chunk)

		// indicate found local chunk
		presentChunks[string(chunkID)] = true

		// sound data chunk is not required if numSampleFrames is zero
		if bytes.Compare(chunkID, COMMONID) == 0 {
			if chunk.(CommonChunk).numSampleFrames == 0 {
				presentChunks[string(SOUNDID)] = true
			}
		}
	}

	// make sure all required chunks are present
	for chunk, present := range presentChunks {
		if !present && utils.ContainsFourCC(RequiredLocalChunkIDs, []byte(chunk)) {
			return FormChunk{}, errors.New(fmt.Sprintf("%s chunk is not present", chunk))
		}
	}

	return form, nil
}
