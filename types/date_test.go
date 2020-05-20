package types

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDate(t *testing.T) {
	const (
		zero  = "0001-01-01"
		value = "2018-01-02"
	)

	var (
		b   []byte
		d   Date
		err error
	)

	t.Run("Parse Date with invalid value", func(it *testing.T) {
		d, err = ParseDate("2020-13-02")
		require.Error(it, err)

		d, err = ParseDate("2020-00-02")
		require.Error(it, err)

		d, err = ParseDate("2020-02-30")
		require.Error(it, err)

		d, err = ParseDate("2020-04-36")
		require.Error(it, err)

		d, err = ParseDate("2020-04--1")
		require.Error(it, err)

		d, err = ParseDate("#$%^&*(")
		require.Error(it, err)
	})

	t.Run("Parse Date with value", func(it *testing.T) {
		d, err := ParseDate("2020-01-01")
		require.NoError(it, err)
		require.Equal(it, Date{Year: 2020, Month: 1, Day: 1}, d)
	})

	t.Run("Date with MarshalJSON", func(it *testing.T) {
		d = Date{}
		require.Equal(it, zero, d.String())
		require.True(it, d.IsZero())

		// marshal with zero value
		b, err = d.MarshalJSON()
		require.NoError(it, err)
		require.Equal(it, strconv.Quote(zero), string(b))
		require.Equal(it, []byte(strconv.Quote(zero)), b)

		d = Date{Year: 2018, Month: 5, Day: 15}
		b, err = d.MarshalJSON()
		require.NoError(it, err)
		require.Equal(it, []byte(`"2018-05-15"`), b)
	})

	t.Run("Date with UnmarshalJSON", func(it *testing.T) {
		d := Date{Year: 20020}

		err = d.UnmarshalJSON([]byte(`null`))
		require.NoError(it, err)
		require.Equal(it, Date{}, d)

		err = d.UnmarshalJSON([]byte(`unknown`))
		require.Error(it, err)

		err = d.UnmarshalJSON([]byte(`"2020-31-12"`))
		require.Error(it, err)

		b = []byte(`"` + value + `"`)
		err = d.UnmarshalJSON(b)
		require.NoError(it, err)
		require.Equal(it, int(2018), d.Year)
		require.Equal(it, int(1), d.Month)
		require.Equal(it, int(2), d.Day)
	})

	t.Run("Date with MarshalText", func(it *testing.T) {
		d = Date{Year: 2018, Month: 1, Day: 2}
		b, err = d.MarshalText()
		require.NoError(it, err)
		require.Equal(it, b, []byte(value))

		d = Date{}
		b, err = d.MarshalText()
		require.NoError(it, err)
		require.Equal(it, b, []byte(zero))
	})

	t.Run("Date with UnmarshalText", func(it *testing.T) {
		b = []byte(value)
		err = d.UnmarshalText(b)
		require.NoError(it, err)
		require.Equal(it, int(2018), d.Year)
		require.Equal(it, int(1), d.Month)
		require.Equal(it, int(2), d.Day)

		err = d.UnmarshalText([]byte(nil))
		require.NoError(it, err)
		require.True(it, d.IsZero())

		err = d.UnmarshalText([]byte(`null`))
		require.NoError(it, err)
		require.True(it, d.IsZero())

		err = d.UnmarshalText([]byte(`03-12-1009`))
		require.Error(it, err)
	})

	t.Run("Date with sql.Scanner", func(it *testing.T) {
		d = Date{}
		dt, _ := time.Parse("2006-01-02", "2018-05-31")
		d.Scan(dt)
		require.Equal(it, Date{Year: 2018, Month: 5, Day: 31}, d)

		d = Date{}
		d.Scan(string("2006-07-26"))
		require.Equal(it, Date{Year: 2006, Month: 7, Day: 26}, d)

		d = Date{}
		d.Scan([]byte("2106-03-26"))
		require.Equal(it, Date{Year: 2106, Month: 3, Day: 26}, d)
	})

	t.Run("Date with driver.Valuer", func(it *testing.T) {
		d = Date{Year: 2018, Month: 5, Day: 31}
		v, err := d.Value()
		require.NoError(it, err)
		require.Equal(it, string("2018-05-31"), v)
	})
}
