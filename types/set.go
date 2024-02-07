package types

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/internal/util"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

type setDataType interface {
	~string
}

// Set : sql data type of `SET`
type Set[T setDataType] []T

var (
	_ db.ColumnDataTyper = (*Set[string])(nil)
)

// DataType :
func (s *Set[T]) ColumnDataType(ctx context.Context) *sql.Column {
	charset, collate := "utf8mb4", "utf8mb4_0900_ai_ci"
	f := sql.GetField(ctx)
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	var def *string
	blr.WriteString("SET(")
	blr.WriteByte('\'')

	val, ok := f.Tag().Option("set")
	if ok {
		paths := strings.Split(val, "|")
		if len(paths) >= 64 {
			panic("maximum 64 of SET value")
		}
		def = &paths[0]
		blr.WriteString(strings.Join(paths, "','"))
	}
	blr.WriteByte('\'')
	blr.WriteByte(')')

	return &sql.Column{
		Name:         f.Name(),
		DataType:     "SET",
		Type:         blr.String(),
		Nullable:     reflext.IsNullable(f.Type()),
		DefaultValue: def,
		Charset:      &charset,
		Collation:    &collate,
	}
}

// Value :
func (s Set[T]) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	for idx, v := range s {
		if idx > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(fmt.Sprintf("%s", v))
	}
	return blr.String(), nil
}

// Scan :
func (s *Set[T]) Scan(it any) error {
	switch vi := it.(type) {
	case []byte:
		s.unmarshal(string(vi))

	case string:
		s.unmarshal(vi)

	case nil:
		*s = nil
	}
	return nil
}

func (s *Set[T]) unmarshal(val string) {
	sets := make(Set[T], 0)
	for _, v := range strings.Split(val, ",") {
		sets = append(sets, T(strings.TrimSpace(v)))
	}
	*s = sets
}
