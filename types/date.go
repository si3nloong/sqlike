package types

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sqlike/columns"
	"github.com/si3nloong/sqlike/util"
)

const dateRegex = `\d{4}\-\d{2}\-\d{2}`

// ErrDateFormat :
var ErrDateFormat = errors.New(`invalid date format, it should be "YYYY-MM-DD"`)

// Date :
type Date struct {
	Year, Month, Day int
}

// DataType :
func (d *Date) DataType(driver string, sf *reflext.StructField) columns.Column {
	dflt := "CURDATE()"
	return columns.Column{
		Name:         sf.Path,
		DataType:     "DATE",
		Type:         "DATE",
		DefaultValue: &dflt,
		Nullable:     sf.IsNullable,
	}
}

// IsZero :
func (d *Date) IsZero() bool {
	return d.Day == 0 && d.Month == 0 && d.Year == 0
}

// Value :
func (d *Date) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan :
func (d *Date) Scan(it interface{}) error {
	switch vi := it.(type) {
	case []byte:
		return d.unmarshal(util.UnsafeString(vi))

	case string:
		return d.unmarshal(vi)

	case time.Time:
		d.Year = vi.Year()
		d.Month = int(vi.Month())
		d.Day = vi.Day()
		return nil

	case nil:
		return nil

	default:
		return errors.New("invalid date format")
	}
}

// String :
func (d *Date) String() string {
	blr := util.AcquireString()
	defer util.ReleaseString(blr)
	d.marshal(blr)
	return blr.String()
}

// MarshalJSON :
func (d *Date) MarshalJSON() ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0, 12))
	b.WriteRune('"')
	d.marshal(b)
	b.WriteRune('"')
	return b.Bytes(), nil
}

// UnmarshalJSON :
func (d *Date) UnmarshalJSON(b []byte) (err error) {
	if b == nil || util.UnsafeString(b) == "null" {
		return nil
	}

	if !regexp.MustCompile(`^\"` + dateRegex + `\"$`).Match(b) {
		return ErrDateFormat
	}

	b = b[1 : len(b)-1]
	return d.unmarshal(util.UnsafeString(b))
}

// MarshalText :
func (d Date) MarshalText() ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0, 10))
	d.marshal(b)
	return b.Bytes(), nil
}

// UnmarshalText :
func (d *Date) UnmarshalText(b []byte) error {
	if b == nil || util.UnsafeString(b) == "null" {
		return nil
	}

	if !regexp.MustCompile(`^` + dateRegex + `$`).Match(b) {
		return ErrDateFormat
	}

	return d.unmarshal(util.UnsafeString(b))
}

func (d *Date) marshal(w writer) {
	year, month, day := 1, 1, 1
	if d.Year > 0 {
		year = d.Year
	}
	if d.Month > 0 {
		month = d.Month
	}
	if d.Day > 0 {
		day = d.Day
	}
	w.WriteString(lpad(strconv.Itoa(year), "0", 4))
	w.WriteByte('-')
	w.WriteString(lpad(strconv.Itoa(month), "0", 2))
	w.WriteByte('-')
	w.WriteString(lpad(strconv.Itoa(day), "0", 2))
}

func (d *Date) unmarshal(str string) (err error) {
	if str == "" {
		return errors.New("types: empty date string")
	}
	// TODO: verify date is valid date
	paths := strings.SplitN(str, "-", 3)
	d.Year, err = strconv.Atoi(paths[0])
	if err != nil {
		return
	}
	d.Month, err = strconv.Atoi(paths[1])
	if err != nil {
		return
	}
	d.Day, err = strconv.Atoi(paths[2])
	if err != nil {
		return
	}
	return
}

func lpad(str, pad string, length int) string {
	for {
		if len(str) >= length {
			return str[0:length]
		}
		str = pad + str
	}
}
