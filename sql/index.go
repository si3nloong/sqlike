package sql

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/bytebufferpool"
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
	MultiValued
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
	case MultiValued:
		return "MULTI-VALUED"
	default:
		return "BTREE"
	}
}

// Index :
type Index struct {
	Name    string
	Cast    string
	As      string
	Type    Type
	Columns []IndexColumn
	Comment string
}

// Direction :
type Direction int

// direction :
const (
	Ascending Direction = iota
	Descending
)

// IndexedColumns :
func IndexedColumns(names ...string) []IndexColumn {
	columns := make([]IndexColumn, 0, len(names))
	for _, n := range names {
		columns = append(columns, IndexedColumn(n))
	}
	return columns
}

// IndexedColumn :
func IndexedColumn(name string) IndexColumn {
	dir := Ascending
	name = strings.TrimSpace(name)
	if name[0] == '-' {
		name = name[1:]
		dir = Descending
	}
	return IndexColumn{
		Name:      name,
		Direction: dir,
	}
}

// IndexColumn :
type IndexColumn struct {
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

func (idx Index) buildName(w interface {
	io.StringWriter
	io.ByteWriter
}) {
	switch idx.Type {
	case Primary:
		w.WriteString("PRIMARY")
		return

	case Unique:
		w.WriteString("UX")

	case FullText:
		w.WriteString("FTX")

	case MultiValued:
		w.WriteString("MVX")

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
	idx.buildName(buf)
	hash.Write(buf.Bytes())
	bytebufferpool.Put(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
