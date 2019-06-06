package types

import (
	"strconv"
	"testing"

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

	t.Run("Date with MarshalJSON", func(it *testing.T) {
		b, err = d.MarshalJSON()
		require.NoError(it, err)
		require.Equal(it, strconv.Quote(zero), string(b))
		require.Equal(it, []byte(strconv.Quote(zero)), b)
	})

	t.Run("Date with UnmarshalJSON", func(it *testing.T) {
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
	})
}
