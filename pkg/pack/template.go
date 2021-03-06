package pack

import (
	"fmt"
	"io"
)

var tplFileReader = `type item struct {
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
}`

var tplGzipReader = `func gzipRead(b []byte) []byte {
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
}`

type TemplateStmt struct {
	Tpl string
}

func (ts *TemplateStmt) Emit(buffer io.Writer) {
	fmt.Fprintf(buffer, "%s\n\n", ts.Tpl)
}
