package pack

import (
	"fmt"
	"io"
)

var tplFileReader = `type file struct {
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
	b, err := ioutil.ReadAll(r.reader)
	if err != nil {
		return make([]byte, 0)
	}
	return b
}

func directReader(b []byte) io.Reader {
	return bytes.NewReader(b)
}`

var tplGzipReader = `var empty = new(emptyBuffer)

type emptyBuffer struct {
}

func (b *emptyBuffer) Read([]byte) (int, error) {
	return 0, io.EOF
}

func gzipReader(b []byte) io.Reader {
	r, err := gzip.NewReader(
		bytes.NewReader(b),
	)
	if err != nil {
		return empty
	}
	defer r.Close()

	uncompressed, err := ioutil.ReadAll(r)
	if err != nil {
		return empty
	}
	return bytes.NewReader(uncompressed)
}`

type TemplateStmt struct {
	Tpl string
}

func (ts *TemplateStmt) Emit(buffer io.Writer) {
	fmt.Fprintf(buffer, "%s\n\n", ts.Tpl)
}
