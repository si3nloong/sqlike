package types

import (
	"database/sql/driver"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/columns"
	"github.com/si3nloong/sqlike/util"
)

type Set []string

// DataType :
func (s Set) DataType(t sqldriver.Info, sf reflext.StructFielder) columns.Column {
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	val, ok := sf.Tag().LookUp("set")
	var def *string
	blr.WriteString("SET(")
	blr.WriteByte('\'')
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
		Charset:      &latin1,
		Collation:    &latin1Bin,
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
		if err := s.unmarshal(string(vi)); err != nil {
			return err
		}
	case string:
		if err := s.unmarshal(vi); err != nil {
			return err
		}
	case nil:
	}
	return nil
}

func (s *Set) unmarshal(val string) error {
	*s = strings.Split(val, ",")
	return nil
}
