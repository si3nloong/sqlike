package types

import (
	"context"
	"database/sql/driver"
	"strings"

	"github.com/si3nloong/sqlike/v2/db"
	"github.com/si3nloong/sqlike/v2/internal/util"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/x/reflext"
)

// Set : sql data type of `SET`
type Set []string

var (
	_ db.ColumnDataTyper = (*Set)(nil)
)

// DataType :
func (s *Set) ColumnDataType(ctx context.Context) *sql.Column {
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
func (s Set) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return strings.Join(s, ","), nil
}

// Scan :
func (s *Set) Scan(it any) error {
	switch vi := it.(type) {
	case []byte:
		s.unmarshal(string(vi))

	case string:
		s.unmarshal(vi)

	case nil:
	}
	return nil
}

func (s *Set) unmarshal(val string) {
	*s = strings.Split(val, ",")
}
