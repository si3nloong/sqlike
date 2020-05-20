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
	Columns []Col
}

// Direction :
type Direction int

// direction :
const (
	Ascending Direction = iota
	Descending
)

// Columns :
func Columns(names ...string) []Col {
	columns := make([]Col, 0, len(names))
	for _, n := range names {
		columns = append(columns, Column(n))
	}
	return columns
}

// Column :
func Column(name string) Col {
	dir := Ascending
	name = strings.TrimSpace(name)
	if name[0] == '-' {
		name = name[1:]
		dir = Descending
	}
	return Col{
		Name:      name,
		Direction: dir,
	}
}

// Col :
type Col struct {
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
