package types

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"

	"github.com/si3nloong/sqlike/reflext"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/columns"
)

const (
	tsFormat      = "2006-01-02 15:04:05"
	zeroTimestamp = "1970-01-01 00:00:01"
)

// Timestamp :
type Timestamp time.Time

var (
	_ driver.Valuer  = (*Timestamp)(nil)
	_ sql.Scanner    = (*Timestamp)(nil)
	_ reflext.Zeroer = (*Timestamp)(nil)
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

// IsZero :
func (ts Timestamp) IsZero() bool {
	t := time.Time(ts)
	if ts.isZero(t) {
		return true
	}
	return t.Format(tsFormat) == zeroTimestamp
}

func (ts Timestamp) isZero(t time.Time) bool {
	return t.Year() <= 1 && t.Month() <= 1 && t.Day() <= 1
}

// String :
func (ts Timestamp) String() string {
	t := time.Time(ts)
	if ts.isZero(t) {
		return zeroTimestamp
	}
	return t.Format(tsFormat)
}

// Value : timestamp alway return UTC
func (ts Timestamp) Value() (driver.Value, error) {
	t := time.Time(ts)
	if t.IsZero() {
		zero, _ := time.Parse(tsFormat, zeroTimestamp)
		return zero, nil
	}
	return t.UTC(), nil
}

// Scan : scan value
func (ts *Timestamp) Scan(it interface{}) error {
	switch vi := it.(type) {
	case []byte:
		t, err := time.Parse(tsFormat, string(vi))
		if err != nil {
			return err
		}
		*ts = Timestamp(t.UTC())

	case string:
		t, err := time.Parse(tsFormat, vi)
		if err != nil {
			return err
		}
		*ts = Timestamp(t.UTC())

	case time.Time:
		*ts = Timestamp(vi.UTC())

	default:
		return errors.New("invalid timestamp format")
	}
	return nil
}
