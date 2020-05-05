package types

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlike/reflext"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/columns"
)

// Timestamp :
type Timestamp time.Time

var (
	_ driver.Valuer = (*Timestamp)(nil)
)

// DataType :
func (ts Timestamp) DataType(_ sqldriver.Info, sf reflext.StructFielder) columns.Column {
	dflt := "NOW()"
	return columns.Column{
		Name:         sf.Name(),
		DataType:     "TIMESTAMP",
		Type:         "TIMESTAMP",
		DefaultValue: &dflt,
		Nullable:     sf.IsNullable(),
	}
}

func (ts Timestamp) Value() (driver.Value, error) {
	t := time.Time(ts)
	if t.IsZero() {
		// zero, _ := time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:01")
		return time.Now().UTC(), nil
	}
	return t, nil
}
