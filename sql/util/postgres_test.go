package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostgresUtil(t *testing.T) {
	utl := PostgresUtil{}

	require.Equal(t, `"abc"`, utl.Quote("abc"))
}
