package types

import (
	"database/sql/driver"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/columns"
	"github.com/si3nloong/sqlike/x/util"
)

// Set : sql data type of `SET`
type Set []string

// DataType :
func (s Set) DataType(_ sqldriver.Info, sf reflext.StructFielder) columns.Column {
	charset, collate := "utf8mb4", "utf8mb4_0900_ai_ci"
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	var def *string
	blr.WriteString("SET(")
	blr.WriteByte('\'')

	val, ok := sf.Tag().LookUp("set")
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

	return columns.Column{
		Name:         sf.Name(),
		Type:         blr.String(),
		DataType:     "SET",
		Nullable:     reflext.IsNullable(sf.Type()),
		Charset:      &charset,
		Collation:    &collate,
		DefaultValue: def,
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
func (s *Set) Scan(it interface{}) error {
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
