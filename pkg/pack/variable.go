package pack

import (
	"fmt"
	"io"
)

type Variable struct {
	Name   string
	Assign string
}

type VariableStmt struct {
	Vars []Variable
}

func (vs *VariableStmt) Emit(buffer io.Writer) {
	if len(vs.Vars) <= 0 {
		return
	}
	vs.begin(buffer)
	for i := 0; i < len(vs.Vars); i++ {
		fmt.Fprintf(buffer, "\t%s = %s\n", vs.Vars[i].Name, vs.Vars[i].Assign)
	}
	vs.end(buffer)
}

func (vs *VariableStmt) begin(buffer io.Writer) {
	fmt.Fprintf(buffer, "var (\n")
}

func (vs *VariableStmt) end(buffer io.Writer) {
	fmt.Fprintf(buffer, ")\n\n")
}
