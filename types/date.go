package types

import (
	"bytes"
	"database/sql/driver"
	"strconv"
	"time"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/sql/component"
	"github.com/si3nloong/sqlike/util"
)

// Date :
type Date struct {
	Year, Month, Day int
}

func (d *Date) init() {
	if d.Year < 1 {
		d.Year = 1
	}
	if d.Month < 1 {
		d.Month = 1
	}
	if d.Day < 1 {
		d.Day = 1
	}
}

// DataType :
func (d *Date) DataType(driver string, sf *reflext.StructField) component.Column {
	dflt := `CURDATE()`
	return component.Column{
		Name:         sf.Path,
		DataType:     `DATE`,
		Type:         `DATE`,
		DefaultValue: &dflt,
		Nullable:     sf.IsNullable,
	}
}

// Value :
func (d *Date) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan :
func (d *Date) Scan(it interface{}) error {
	// TODO:
	switch vi := it.(type) {
	case []byte:

	case string:
	case time.Time:
		d.Year = vi.Year()
		d.Month = int(vi.Month())
		d.Day = vi.Day()
	}
	return nil
}

// String :
func (d *Date) String() string {
	d.init()
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	blr.WriteString(lpad(strconv.Itoa(d.Year), "0", 4))
	blr.WriteRune('-')
	blr.WriteString(lpad(strconv.Itoa(d.Month), "0", 2))
	blr.WriteRune('-')
	blr.WriteString(lpad(strconv.Itoa(d.Day), "0", 2))
	return blr.String()
}

// MarshalJSON :
func (d *Date) MarshalJSON() ([]byte, error) {
	d.init()
	b := bytes.NewBuffer(make([]byte, 0, 12))
	b.WriteRune('"')
	b.WriteString(lpad(strconv.Itoa(d.Year), "0", 4))
	b.WriteRune('-')
	b.WriteString(lpad(strconv.Itoa(d.Month), "0", 2))
	b.WriteRune('-')
	b.WriteString(lpad(strconv.Itoa(d.Day), "0", 2))
	b.WriteRune('"')
	return b.Bytes(), nil
}

// UnmarshalJSON :
func (d *Date) UnmarshalJSON(b []byte) error {
	return nil
}

func lpad(str, pad string, length int) string {
	for {
		str = pad + str
		if len(str) >= length {
			return str[0:length]
		}
	}
}
