package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoolean(t *testing.T) {
	var (
		flag Boolean
		err  error
	)

	t.Run("Scan with []byte", func(it *testing.T) {
		err = flag.Scan([]byte(`Yes`))
		require.NoError(t, err)
		require.Equal(t, Boolean(true), flag)

		err = flag.Scan([]byte(`yEs`))
		require.NoError(t, err)
		require.Equal(t, Boolean(true), flag)

		err = flag.Scan([]byte(`YES`))
		require.NoError(t, err)
		require.Equal(t, Boolean(true), flag)

		err = flag.Scan([]byte(`y`))
		require.NoError(t, err)
		require.Equal(t, Boolean(true), flag)

		err = flag.Scan([]byte(`true`))
		require.NoError(t, err)
		require.Equal(t, Boolean(true), flag)

		err = flag.Scan([]byte(`no`))
		require.NoError(t, err)
		require.Equal(t, Boolean(false), flag)

		err = flag.Scan([]byte(`n`))
		require.NoError(t, err)
		require.Equal(t, Boolean(false), flag)

		err = flag.Scan([]byte(`false`))
		require.NoError(t, err)
		require.Equal(t, Boolean(false), flag)
	})

	t.Run("Scan with string", func(it *testing.T) {
		err = flag.Scan("yes")
		require.NoError(t, err)
		require.Equal(t, Boolean(true), flag)

		err = flag.Scan("y")
		require.NoError(t, err)
		require.Equal(t, Boolean(true), flag)

		err = flag.Scan("true")
		require.NoError(t, err)
		require.Equal(t, Boolean(true), flag)

		err = flag.Scan("no")
		require.NoError(t, err)
		require.Equal(t, Boolean(false), flag)

		err = flag.Scan("n")
		require.NoError(t, err)
		require.Equal(t, Boolean(false), flag)

		err = flag.Scan("false")
		require.NoError(t, err)
		require.Equal(t, Boolean(false), flag)
	})

	t.Run("Scan with int64", func(it *testing.T) {
		err = flag.Scan(int64(1))
		require.NoError(t, err)
		require.Equal(t, Boolean(true), flag)

		err = flag.Scan(int64(0))
		require.NoError(t, err)
		require.Equal(t, Boolean(false), flag)
	})
}
