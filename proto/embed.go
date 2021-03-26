package proto

import (
	"compress/gzip"
	"io"
	"io/ioutil"
)

var (
	bu = []byte{}
)

var (
	Config = file{"conf/session.yaml", gzipReader(&buffer{data: bu, index: 0})}
)

type buffer struct {
	data  []byte
	index int64
}

func (b *buffer) Read(p []byte) (int, error) {
	if len(b.data) <= 0 {
		return 0, io.EOF
	}
	if b.index >= int64(len(b.data)) {
		return 0, io.EOF
	}
	// Copy
	n := copy(p, b.data[b.index:])
	b.index += int64(n)

	return n, nil
}

type emptyBuffer struct {
}

func (b *emptyBuffer) Read([]byte) (int, error) {
	return 0, io.EOF
}

var invalidBuffer = new(emptyBuffer)

func gzipReader(in io.Reader) io.Reader {
	r, err := gzip.NewReader(in)
	if err != nil {
		return invalidBuffer
	}
	return r
}

type file struct {
	file   string
	reader io.Reader
}

func (r *file) Read(buffer []byte) (int, error) {
	return r.reader.Read(buffer)
}

func (r *file) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(r.reader)
}
