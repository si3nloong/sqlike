package types

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimestamp(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	require.NoError(t, err)

	utcNow := time.Now().UTC()

	t.Run("DataType", func(it *testing.T) {
		var ts Timestamp

		col := ts.DataType(nil, field{
			name: "Timestamp",
			t:    reflect.TypeOf(ts),
		})

		require.Equal(it, "Timestamp", col.Name)
		require.Equal(it, "TIMESTAMP", col.DataType)
		require.Equal(it, "TIMESTAMP", col.Type)
		require.Equal(it, "NOW()", *col.DefaultValue)
		require.False(it, col.Nullable)
		require.Empty(it, col.Charset)
		require.Empty(it, col.Collation)
	})

	t.Run("IsZero", func(it *testing.T) {
		var ts Timestamp

		require.True(it, ts.IsZero())
		require.Equal(it, "1970-01-01 00:00:01", ts.String())
	})

	t.Run("driver.Valuer", func(it *testing.T) {
		ts := Timestamp(utcNow.In(loc))
		v, err := ts.Value()
		require.NoError(t, err)
		require.Equal(t, utcNow.Format(time.RFC3339), v.(time.Time).Format(time.RFC3339))
	})

	t.Run("sql.Scanner", func(it *testing.T) {
		// scan with []byte
		{
			var ts Timestamp
			dt, _ := time.Parse("2006-01-02 15:04:05", "2020-05-03 16:00:00")
			err := ts.Scan([]byte("2020-05-03 16:00:00"))
			require.NoError(it, err)
			require.Equal(it, Timestamp(dt), ts)
		}

		// scan with string
		{
			var ts Timestamp
			dt, _ := time.Parse("2006-01-02 15:04:05", "2020-05-03 16:00:00")
			err := ts.Scan("2020-05-03 16:00:00")
			require.NoError(it, err)
			require.Equal(it, Timestamp(dt), ts)
		}

		// scan with time.Time
		{
			var ts Timestamp
			dt, _ := time.Parse("2006-01-02", "2018-09-04")
			err := ts.Scan(dt)
			require.NoError(it, err)
			require.Equal(it, Timestamp(dt), ts)
		}

		// scan with error
		{
			var ts Timestamp
			err := ts.Scan(true)
			require.Error(it, err)
		}
	})

}
