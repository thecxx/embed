package pack

import (
	"fmt"
	"io"
)

var tplBuffer = `type buffer struct {
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

func (r *file) Bytes() []byte {
	if b, err := ioutil.ReadAll(r.reader); err == nil {
		return b
	}
	return make([]byte, 0)
}`

var tplGzipReader = `type emptyBuffer struct {
}

func (b *emptyBuffer) Read([]byte) (int, error) {
	return 0, io.EOF
}

func gzipReader(in io.Reader) io.Reader {
	r, err := gzip.NewReader(in)
	if err != nil {
		return new(emptyBuffer)
	}
	return r
}`

type TemplateStmt struct {
	Tpl string
}

func (ts *TemplateStmt) Emit(buffer io.Writer) {
	fmt.Fprintf(buffer, "%s\n\n", ts.Tpl)
}
