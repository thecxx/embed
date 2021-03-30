package pack

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

type Statement interface {
	Emit(buffer io.Writer)
}

type SourceFile struct {
	deps   []string
	stmts  []Statement
	buffer *bytes.Buffer
}

// NewSourceFile returns a source file object.
func NewSourceFile(pkg string) *SourceFile {
	return &SourceFile{
		buffer: bytes.NewBufferString(fmt.Sprintf("package %s\n\n", pkg)),
	}
}

// Import imports a package.
func (sf *SourceFile) Import(path string) {
	for i := 0; i < len(sf.deps); i++ {
		if sf.deps[i] == path {
			return
		}
	}
	sf.deps = append(sf.deps, path)
}

func (sf *SourceFile) DeclareVar(vars ...Variable) {
	sf.AddStmt(&VariableStmt{vars})
}

func (sf *SourceFile) DeclareGzipReader() {
	sf.Import("compress/gzip")
	sf.AddStmt(&TemplateStmt{tplGzipReader})
}

func (sf *SourceFile) DeclareFileReader() {
	sf.AddStmt(&TemplateStmt{tplFileReader})
}

func (sf *SourceFile) Bytes() []byte {
	sf.build()
	return sf.buffer.Bytes()
}

func (sf *SourceFile) String() string {
	sf.build()
	return sf.buffer.String()
}

func (sf *SourceFile) build() {
	// Statement: import
	if len(sf.deps) > 0 {
		sort.Strings(sf.deps)
		// import ()
		fmt.Fprintf(sf.buffer, "import (\n")
		for i := 0; i < len(sf.deps); i++ {
			fmt.Fprintf(sf.buffer, "\t\"%s\"\n", sf.deps[i])
		}
		fmt.Fprintf(sf.buffer, ")\n\n")
	}
	// Other statements
	for _, stmt := range sf.stmts {
		stmt.Emit(sf.buffer)
	}
}

func (sf *SourceFile) AddStmt(stmt Statement) {
	sf.stmts = append(sf.stmts, stmt)
}
