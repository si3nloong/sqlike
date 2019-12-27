package indexes

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/bytebufferpool"
)

type writer interface {
	io.Writer
	io.StringWriter
	io.ByteWriter
}

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
func Columns(cols ...string) []Column {
	columns := make([]Column, 0, len(cols))
	for _, col := range cols {
		dir := Ascending
		col = strings.TrimSpace(col)
		if col[0] == '-' {
			col = col[1:]
			dir = Descending
		}
		columns = append(columns, Column{
			Name:      col,
			Direction: dir,
		})
	}
	return columns
}

// Column :
type Column struct {
	Name      string
	Direction Direction
}

// GetName :
func (idx Index) GetName() string {
	if idx.Name != "" {
		return idx.Name
	}
	return idx.HashName()
}

func (idx Index) buildName(w writer) {
	switch idx.Type {
	case Unique:
		w.WriteString("UX")
	case Primary:
		w.WriteString("PRIMARY")
		return
	default:
		w.WriteString("IX")
	}
	w.WriteByte('-')
	for i, col := range idx.Columns {
		if i > 0 {
			w.WriteByte(';')
		}
		w.WriteString(col.Name)
		w.WriteByte('@')
		if col.Direction == 0 {
			w.WriteString("ASC")
		} else {
			w.WriteString("DESC")
		}
	}
}

// HashName :
func (idx Index) HashName() string {
	hash := md5.New()
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	idx.buildName(buf)
	hash.Write(buf.Bytes())
	return fmt.Sprintf("%x", hash.Sum(nil))
}
