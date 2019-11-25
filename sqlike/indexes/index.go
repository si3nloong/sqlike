package indexes

import (
	"github.com/si3nloong/sqlike/util"
)

// Type :
type Type int

// types :
const (
	BTree Type = iota + 1
	FullText
	Unique
	Spatial
	Primary
)

func (t Type) String() string {
	switch t {
	case FullText:
		return "FULLTEXT"
	case Unique:
		return "UNIQUE"
	case Spatial:
		return "SPATIAL"
	case Primary:
		return "PRIMARY"
	default:
		return "BTREE"
	}
}

// Index :
type Index struct {
	Name    string
	Type    Type
	Columns []Column
}

// Direction :
type Direction int

// direction :
const (
	Ascending Direction = iota
	Descending
)

// Columns :
func Columns(cols ...interface{}) []Column {
	return nil
}

// Column :
type Column struct {
	Name      string
	Direction Direction
}

// GetName :
func (idx *Index) GetName() string {
	if idx.Name != "" {
		return idx.Name
	}
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	switch idx.Type {
	case Unique:
		blr.WriteString("UX")
	case Primary:
		blr.WriteString("PRIMARY")
		return ""
	default:
		blr.WriteString("IX")
	}
	blr.WriteByte('-')
	for i, col := range idx.Columns {
		if i > 0 {
			blr.WriteByte('-')
		}
		blr.WriteString(col.Name)
		blr.WriteByte('_')
		if col.Direction == 0 {
			blr.WriteString("ASC")
		} else {
			blr.WriteString("DESC")
		}
	}
	idx.Name = blr.String()
	return idx.Name
}
