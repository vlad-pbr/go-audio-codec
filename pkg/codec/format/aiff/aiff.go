package aiff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec/audio"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/utils"
)

var FORMID utils.FourCC = []byte("FORM")
var COMMONID utils.FourCC = []byte("COMM")
var SOUNDID utils.FourCC = []byte("SSND")

var FORMTYPE utils.FourCC = []byte("AIFF")

// TODO float80
type extended []byte

type AIFFFormat struct {
}

type formChunk struct {
	utils.Chunk
	formType    utils.FourCC // must be AIFF
	localChunks []utils.ChunkInterface
}

func newFormChunk(buffer *bytes.Buffer) (formChunk, error) {

	var form formChunk

	// parse form chunk ID
	form.ChunkID = buffer.Next(4)
	if bytes.Compare(form.ChunkID, FORMID) != 0 {
		return formChunk{}, errors.New(fmt.Sprintf("FORM chunk ID is invalid: found %s, must be %s", form.ChunkID, FORMID))
	}

	// parse form chunk size
	form.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))

	// parse form chunk ID
	form.formType = buffer.Next(4)
	if bytes.Compare(form.formType, FORMTYPE) != 0 {
		return formChunk{}, errors.New(fmt.Sprintf("FORM chunk type is invalid: found %s, must be %s", form.formType, FORMTYPE))
	}

	// the following chunks must be present
	var requiredChunks = map[string]bool{
		string(COMMONID): false,
		string(SOUNDID):  false,
	}

	// while unread portion of buffer is not empty
	for buffer.Len() != 0 {

		chunkID := buffer.Next(4)

		// if curren chunk is COMMON
		if bytes.Compare(chunkID, COMMONID) == 0 {

			// make sure not already found
			if requiredChunks[string(COMMONID)] {
				return formChunk{}, errors.New("more than one instance of COMMON chunk present")
			}

			// create new common chunk
			commonChunk, err := newCommonChunk(buffer)
			if err != nil {
				return formChunk{}, errors.New(fmt.Sprintf("error occurred while decoding COMMON chunk: %s", err.Error()))
			}

			// append resulting common chunk
			form.localChunks = append(form.localChunks, commonChunk)

			// indicate found common chunk
			requiredChunks[string(COMMONID)] = true

			// sound data chunk is not required if numSampleFrames is zero
			if commonChunk.numSampleFrames == 0 {
				requiredChunks[string(SOUNDID)] = true
			}

		} else if bytes.Compare(chunkID, SOUNDID) == 0 {

			// make sure not already found
			if requiredChunks[string(SOUNDID)] {
				return formChunk{}, errors.New("more than one instance of SOUND chunk present")
			}

			// create new sound chunk
			soundDataChunk, err := newSoundChunk(buffer)
			if err != nil {
				return formChunk{}, errors.New(fmt.Sprintf("error occurred while decoding SOUND chunk: %s", err.Error()))
			}

			// append resulting sound chunk
			form.localChunks = append(form.localChunks, soundDataChunk)

			// indicate found sound chunk
			requiredChunks[string(SOUNDID)] = true

		} else {
			// skip current chunk
			buffer.Next(int(int32(binary.BigEndian.Uint32(buffer.Next(4)))))
		}

	}

	// make sure all required chunks are present
	for chunk, present := range requiredChunks {
		if !present {
			return formChunk{}, errors.New(fmt.Sprintf("%s chunk is not present", chunk))
		}
	}

	return form, nil
}

// TODO multichannel table

type commonChunk struct {
	utils.Chunk
	numChannels     int16
	numSampleFrames uint32
	sampleSize      int16
	sampleRate      extended
}

func (c commonChunk) GetBytes() []byte {
	return c.GetBytesWithHeaders(c.numChannels, c.numSampleFrames, c.sampleSize, c.sampleRate)
}

func newCommonChunk(buffer *bytes.Buffer) (commonChunk, error) {

	// define common chunk struct
	var commChunk commonChunk
	commChunk.ChunkID = COMMONID
	commChunk.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))

	// make sure common chunk size is 18
	if commChunk.ChunkSize != 18 {
		return commonChunk{}, errors.New(fmt.Sprintf("COMMON chunk size is invalid: found %d, must be %d", commChunk.ChunkSize, 18))
	}

	// fill common chunk struct
	commChunk.numChannels = int16(binary.BigEndian.Uint16(buffer.Next(2)))
	commChunk.numSampleFrames = binary.BigEndian.Uint32(buffer.Next(4))
	commChunk.sampleSize = int16(binary.BigEndian.Uint16(buffer.Next(2)))
	commChunk.sampleRate = buffer.Next(10)

	return commChunk, nil
}

type soundDataChunk struct {
	utils.Chunk
	offset    uint32
	blockSize uint32
	soundData []uint8
}

func (c soundDataChunk) GetBytes() []byte {
	return c.GetBytesWithHeaders(c.offset, c.blockSize, c.soundData)
}

func newSoundChunk(buffer *bytes.Buffer) (soundDataChunk, error) {

	// fill common chunk struct
	var soundChunk soundDataChunk
	soundChunk.ChunkID = SOUNDID
	soundChunk.ChunkSize = int32(binary.BigEndian.Uint32(buffer.Next(4)))
	soundChunk.offset = binary.BigEndian.Uint32(buffer.Next(4))
	soundChunk.blockSize = binary.BigEndian.Uint32(buffer.Next(4))

	// parse sound chunk samples
	for i := 8; i != int(soundChunk.ChunkSize); i++ {
		sample, err := buffer.ReadByte()
		if err != nil {
			return soundDataChunk{}, errors.New(fmt.Sprintf("unexpected EOF while reading SOUND chunk samples"))
		}
		soundChunk.soundData = append(soundChunk.soundData, uint8(sample))
	}

	return soundChunk, nil
}

// TODO optional chunks

func (f AIFFFormat) Decode(data []byte) (audio.Audio, error) {

	// create form chunk
	formChunk, err := newFormChunk(bytes.NewBuffer(data))
	if err != nil {
		return audio.Audio{}, err
	}

	// define audio struct
	audio := audio.Audio{}

	// iterate form local chunks and fill audio struct accordingly
	for _, chunk := range formChunk.localChunks {
		if bytes.Compare(chunk.GetID(), COMMONID) == 0 {
			audio.NumChannels = uint16(chunk.(commonChunk).numChannels)
		}
	}

	return audio, nil
}
