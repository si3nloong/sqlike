package rsql

// Ast :
type Ast struct {
	Token  string
	Nodes  []*Ast
	Parent *Ast
}
