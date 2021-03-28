package stmt

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type dependent struct {
	pkg  string
	path string
}

type TopStatement struct {
	pkg   string
	stmts []Statement
	deps  map[string]dependent
}

// NewTopStatement returns a top statement.
func NewTopStatement(pkg string) *TopStatement {
	return &TopStatement{
		pkg:   pkg,
		stmts: make([]Statement, 0),
		deps:  make(map[string]dependent),
	}
}

func (ts *TopStatement) Import(path, pkg string) {
	if _, ok := ts.deps[path]; !ok {
		ts.deps[path] = dependent{pkg, path}
	}
}

func (ts *TopStatement) Push(stmt Statement) {
	ts.stmts = append(ts.stmts, stmt)
}

func (ts *TopStatement) Export() string {
	// package {name}
	page := bytes.NewBufferString(
		fmt.Sprintf("package %s\n\n", ts.pkg),
	)

	// import (
	//     {name} {path}
	// )
	if len(ts.deps) > 0 {
		imports := make([]string, 0)
		for _, dep := range ts.deps {
			if len(dep.pkg) <= 0 {
				imports = append(imports, dep.path)
			} else {
				imports = append(imports, dep.pkg+" "+dep.path)
			}
		}
		sort.Strings(imports)
		fmt.Fprintf(page, "import (\n")
		for _, i := range imports {
			if strings.Contains(i, " ") {
				strs := strings.SplitN(i, " ", 2)
				fmt.Fprintf(page, "    %s \"%s\"\n", strs[0], strs[1])
			} else {
				fmt.Fprintf(page, "    \"%s\"\n", i)
			}
		}
		fmt.Fprintf(page, ")\n\n")
	}

	// Other statements
	for _, stmt := range ts.stmts {
		fmt.Fprintf(page, "%s\n", stmt.Block())
	}

	// Buffer support
	fmt.Fprintf(page, `type buffer struct {
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
}
`)

	return page.String()
}

type Statement interface {
	Block() string
}
