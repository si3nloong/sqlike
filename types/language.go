package types

import (
	"database/sql/driver"
	"strings"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/columns"
	"golang.org/x/text/language"
)

// Language :
type Language struct {
	language.Tag
}

// DataType :
func (l Language) DataType(driver string, sf *reflext.StructField) columns.Column {
	return columns.Column{
		Name:      sf.Path,
		DataType:  "CHAR",
		Type:      "CHAR(3)",
		Nullable:  reflext.IsNullable(sf.Zero.Type()),
		CharSet:   &latin1,
		Collation: &latin1Bin,
	}
}

// Value :
func (l Language) Value() (driver.Value, error) {
	return strings.ToUpper(l.String()), nil
}

// Scan :
func (l *Language) Scan(it interface{}) error {
	switch vi := it.(type) {
	case []byte:
		lg, err := language.Parse(string(vi))
		if err != nil {
			return err
		}
		*l = Language{lg}
	case string:
		lg, err := language.Parse(vi)
		if err != nil {
			return err
		}
		*l = Language{lg}
	default:
	}
	return nil
}
