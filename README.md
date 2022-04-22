# GO Audio Codec

Audio codec written purely in go. Currently supports encoding and decoding for the following audio formats:
- WAV (Canonical)
- AIFF

This is a personal project developed as a part of a bigger project of mine. Documentation for different parts of the codec is available in READMEs within subdirectories as well as [docs](/docs).

## Example

Here is how you can use this codec to convert an AIFF file to a WAV file:

``` golang
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/vlad-pbr/go-audio-codec/pkg/codec"
	"github.com/vlad-pbr/go-audio-codec/pkg/codec/format"
)

func main() {

    // read audio file as bytes
    data, err := ioutil.ReadFile("path/to/file.aif")
    if err != nil {
        panic(fmt.Errorf("could not read file: %s", err.Error()))
    }

    // decode audio bytes to an audio container
    aud, err := codec.Decode(data)
    if err != nil {
        panic(fmt.Errorf("could not decode audio: %s", err.Error()))
    }

    // display audio info
    fmt.Println(aud)

    // write audio in WAV format
    out := codec.Encode(aud, format.WAV)
    if err := ioutil.WriteFile("path/to/file.wav", out, 0664); err != nil {
    	panic(fmt.Errorf("could not write file: %s", err.Error()))
    }

}
```