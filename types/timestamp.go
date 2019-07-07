package types

import (
	"time"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/columns"
)

// Timestamp :
type Timestamp time.Time

// DataType :
func (ts *Timestamp) DataType(driver string, sf *reflext.StructField) columns.Column {
	dflt := "CURDATE()"
	return columns.Column{
		Name:         sf.Path,
		DataType:     "TIMESTAMP",
		Type:         "TIMESTAMP",
		DefaultValue: &dflt,
		Nullable:     sf.IsNullable,
	}
}
