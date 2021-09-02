package aiff

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var FORMID utils.FourCC = [4]byte{70, 79, 82, 77}   // FORM
var FORMTYPE utils.FourCC = [4]byte{65, 73, 70, 70} // AIFF

// possible local chunks, in order of precedence
var LocalChunks = map[string]interface{}{
	string(COMMONID[:]):              NewCommonChunk,
	string(SOUNDID[:]):               NewSoundChunk,
	string(MARKERID[:]):              NewMarkerChunk,
	string(INSTRUMENTID[:]):          NewInstrumentChunk,
	string(COMMENTSID[:]):            NewCommentsChunk,
	string(NAMEID[:]):                NewNameChunk,
	string(AUTHORID[:]):              NewAuthorChunk,
	string(COPYRIGHTID[:]):           NewCopyrightChunk,
	string(ANNOTATIONID[:]):          NewAnnotationChunk,
	string(AUDIORECORDINGID[:]):      NewAudioRecordingChunk,
	string(MIDIDATAID[:]):            NewMIDIDataChunk,
	string(APPLICATIONSPECIFICID[:]): NewApplicationSpecificChunk,
}

// local chunks which are required in a valid AIFF file
var RequiredLocalChunkIDs = []utils.FourCC{
	COMMONID,
	SOUNDID,
}

// chunks which are allowed to be present more than once
var AllowedMultipleChunks = []utils.FourCC{
	ANNOTATIONID,
	MIDIDATAID,
	APPLICATIONSPECIFICID,
}

type FormChunk struct {
	AIFFChunk
	FormType    utils.FourCC // must be AIFF
	LocalChunks []utils.ChunkInterface
}

func (c FormChunk) GetBytes() []byte {
	return c.MakeChunkBytes(
		c.FormType,
		utils.GetChunksBytes(c.LocalChunks),
	)
}

func FormHeaders(buffer *bytes.Buffer) (utils.FourCC, int32, utils.FourCC, error) {

	var chunkID utils.FourCC
	var chunkSize int32
	var formType utils.FourCC

	// parse form chunk ID
	copy(chunkID[:], utils.Next(buffer, 4))
	if !bytes.Equal(chunkID[:], FORMID[:]) {
		return chunkID, chunkSize, formType, fmt.Errorf("FORM chunk ID is invalid: found %s, must be %s", chunkID, FORMID)
	}

	// parse form chunk size
	chunkSize = int32(binary.BigEndian.Uint32(utils.Next(buffer, 4)))

	// parse form type
	copy(formType[:], utils.Next(buffer, 4))
	if !bytes.Equal(formType[:], FORMTYPE[:]) {
		return chunkID, chunkSize, formType, fmt.Errorf("FORM chunk type is invalid: found %s, must be %s", formType, FORMTYPE)
	}

	return chunkID, chunkSize, formType, nil
}

func NewFormChunk(buffer *bytes.Buffer) (FormChunk, error) {

	var form FormChunk
	var err error

	// parse form chunk headers
	form.ChunkID, form.ChunkSize, form.FormType, err = FormHeaders(buffer)
	if err != nil {
		return form, fmt.Errorf("error while decoding FORM chunk headers: %s", err.Error())
	}

	// the following chunks can be present
	var presentChunks = map[string]bool{}
	for fourCC := range LocalChunks {
		presentChunks[fourCC] = false
	}

	// read until end of buffer (account for zero padding)
	for buffer.Len() > 1 {

		var chunkID utils.FourCC
		copy(chunkID[:], utils.Next(buffer, 4))

		// make sure chunk not already present unless allowed to be present multiple times
		if presentChunks[string(chunkID[:])] && !utils.ContainsFourCC(AllowedMultipleChunks, chunkID) {
			return form, fmt.Errorf("more than one instance of %s local chunk present", string(chunkID[:]))
		}

		// retrieve target chunk's creation function
		newChunkFunction := LocalChunks[string(chunkID[:])]
		if newChunkFunction == nil {
			return form, fmt.Errorf("invalid local chunk ID found: %s", string(chunkID[:]))
		}

		// create target chunk
		chunk, err := newChunkFunction.(func(*bytes.Buffer) (utils.ChunkInterface, error))(buffer)
		if err != nil {
			return form, fmt.Errorf("error occurred while decoding %s local chunk: %s", string(chunkID[:]), err.Error())
		}

		// append resulting local chunk
		form.LocalChunks = append(form.LocalChunks, chunk)

		// indicate found local chunk
		presentChunks[string(chunkID[:])] = true

		// sound data chunk is not required if numSampleFrames is zero
		if bytes.Equal(chunkID[:], COMMONID[:]) {
			if chunk.(CommonChunk).NumSampleFrames == 0 {
				presentChunks[string(SOUNDID[:])] = true
			}
		}
	}

	adjustForZeroPadding(form.ChunkSize, buffer)

	// make sure all required chunks are present
	for chunk, present := range presentChunks {
		var chunkFourCC utils.FourCC
		copy(chunkFourCC[:], chunk[:])

		if !present && utils.ContainsFourCC(RequiredLocalChunkIDs, chunkFourCC) {
			return form, fmt.Errorf("%s local chunk is not present", chunk)
		}
	}

	return form, nil
}
