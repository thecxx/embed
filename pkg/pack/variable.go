package pack

import (
	"fmt"
	"io"
)

type Variable struct {
	Name    string
	Assign  string
	Comment string
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
		v := vs.Vars[i]
		if len(v.Comment) > 0 {
			fmt.Fprintf(buffer, "\t// %s\n", v.Comment)
		}
		fmt.Fprintf(buffer, "\t%s = %s\n", v.Name, v.Assign)
	}
	vs.end(buffer)
}

func (vs *VariableStmt) begin(buffer io.Writer) {
	fmt.Fprintf(buffer, "var (\n")
}

func (vs *VariableStmt) end(buffer io.Writer) {
	fmt.Fprintf(buffer, ")\n\n")
}
