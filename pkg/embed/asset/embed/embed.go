// Time: 2021-03-31 16:07:16 +0800 CST
// Generated by embed, DO NOT EDIT.

package embed

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
)

type item struct {
	file string
	data []byte
}

// NewReader returns a new reader for the file.
func (i *item) NewReader() io.Reader {
	return bytes.NewReader(i.data)
}

// File returns the file path.
func (i *item) File() string {
	return i.file
}

// Size returns the size.
func (i *item) Size() int {
	return len(i.data)
}

// Bytes returns the data.
func (i *item) Bytes() []byte {
	return i.data
}

func directRead(b []byte) []byte {
	if len(b) <= 0 {
		return make([]byte, 0)
	}
	return b
}

func gzipRead(b []byte) []byte {
	r, err := gzip.NewReader(
		bytes.NewReader(b),
	)
	if err != nil {
		return make([]byte, 0)
	}
	defer r.Close()

	uncompressed, err := ioutil.ReadAll(r)
	if err != nil {
		return make([]byte, 0)
	}

	return uncompressed
}


// An embedded resource configuration.
var EmbedYaml = item{"embed.yaml", gzipRead(_5b6e04661e88337a42d767ea6966bcdb)}
