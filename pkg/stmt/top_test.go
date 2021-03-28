package stmt

import (
	"fmt"
	"testing"
)

func TestTopStatement_Export(t *testing.T) {
	ts := NewTopStatement("test")
	ts.Import("github.com/test/test1", "")
	ts.Import("github.com/test/test2", "t2")

	fmt.Printf("%s\n", ts.Export())
}
