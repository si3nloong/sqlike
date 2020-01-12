package examples

import (
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/stretchr/testify/require"
)

// DatabaseExamples :
func DatabaseExamples(t *testing.T, client *sqlike.Client) {
	var (
		err error
	)

	{
		err = client.DropDatabase("test")
		require.NoError(t, err)
		err = client.CreateDatabase("test")
		require.NoError(t, err)

		testDB := client.Database("test")
		tb := testDB.Table("t1")
		require.NotNil(t, tb)
	}

	{
		dbs, err := client.ListDatabases()
		require.True(t, len(dbs) > 0)
		require.NoError(t, err)
	}
}
