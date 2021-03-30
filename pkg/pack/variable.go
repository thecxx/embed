package pack

import (
	"fmt"
	"io"
	"strings"
)

type Variable struct {
	Name    string
	Comment string
	Assign  string
}

type VariableStmt struct {
	Vars []Variable
}

func (vs *VariableStmt) Emit(buffer io.Writer) {
	if len(vs.Vars) <= 0 {
		return
	}
	// Just only one
	if len(vs.Vars) == 1 {
		if len(vs.Vars[0].Comment) > 0 {
			fmt.Fprintf(buffer, "\n// %s\n", strings.Replace(vs.Vars[0].Comment, "\n", "\n// ", -1))
		}
		fmt.Fprintf(buffer, "var %s = %s\n", vs.Vars[0].Name, vs.Vars[0].Assign)
	} else {
		vs.begin(buffer)
		for i := 0; i < len(vs.Vars); i++ {
			v := vs.Vars[i]
			if len(v.Comment) > 0 {
				fmt.Fprintf(buffer, "\n\t// %s\n", strings.Replace(v.Comment, "\n", "\n\t// ", -1))
			}
			fmt.Fprintf(buffer, "\t%s = %s\n", v.Name, v.Assign)
		}
		vs.end(buffer)
	}
}

func (vs *VariableStmt) begin(buffer io.Writer) {
	fmt.Fprintf(buffer, "var (\n")
}

func (vs *VariableStmt) end(buffer io.Writer) {
	fmt.Fprintf(buffer, ")\n\n")
}
