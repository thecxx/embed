package stmt

type vardecl struct {
	name   string
	ty     string
	assign string
}

type VarStatement struct {
	vars []vardecl
}

func (vs *VarStatement) Block() string {
	return ""
}
