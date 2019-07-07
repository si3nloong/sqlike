package indexes

import "strings"

// Kind :
type Kind int

// types :
const (
	FullText Kind = iota + 1
	Unique
	Spatial
)

// Order :
type Order int

// ordering
const (
	Desc Order = iota
	Asc
)

// Index :
type Index struct {
	Name    string
	Kind    Kind
	Order   Order
	Columns []string
}

// GetName :
func (idx *Index) GetName() string {
	if idx.Name == "" {
		name := strings.Join(idx.Columns, "_")
		switch idx.Kind {
		case FullText:
		case Unique:
			idx.Name = "UX_" + name
		default:
			idx.Name = "IX_" + name
		}
	}
	return idx.Name
}
