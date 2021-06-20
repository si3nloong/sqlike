package examples

import (
	"context"
	"testing"

	"github.com/si3nloong/sqlike"
	"github.com/stretchr/testify/require"
)

// DatabaseExamples :
func DatabaseExamples(t *testing.T, client *sqlike.Client) {
	var (
		err error
		ctx = context.Background()
	)

	{
		err = client.DropDatabase(ctx, "test")
		require.NoError(t, err)
		err = client.CreateDatabase(ctx, "test")
		require.NoError(t, err)

		testDB := client.Database("test")
		tb := testDB.Table("t1")
		require.NotNil(t, tb)
	}

	{
		dbs, err := client.ListDatabases(ctx)
		require.True(t, len(dbs) > 0)
		require.NoError(t, err)
	}
}
