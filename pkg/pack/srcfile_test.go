package pack

import (
	"fmt"
	"testing"
)

func TestSourceFile_String(t *testing.T) {
	sf := NewSourceFile("test")

	sf.Import("os")
	sf.Import("compress/gzip")

	sf.DeclareVar(
		Variable{Name: "TestFile", Assign: "file{\"test.txt\", &buffer{data: xxxxx, index: 0}}"},
		Variable{Name: "TestFile1", Assign: "file{\"test.txt\", &buffer{data: xxxxx, index: 0}}"},
	)

	sf.DeclareGzipReader()

	fmt.Printf("%s\n", sf.String())

}
