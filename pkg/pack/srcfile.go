package pack

import (
	"bytes"
	"fmt"
)

type SourceFile struct {
	buffer *bytes.Buffer
}

func NewSourceFile(pkg string) *SourceFile {
	return &SourceFile{
		buffer: bytes.NewBufferString(fmt.Sprintf("package %s\n\n", pkg)),
	}
}

func (sf *SourceFile) Bytes() []byte {
	return sf.buffer.Bytes()
}

func (sf *SourceFile) String() string {
	return sf.buffer.String()
}
