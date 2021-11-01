package examples

import (
	"context"
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/stretchr/testify/require"
)

// DatabaseExamples :
func DatabaseExamples(t *testing.T, client *sqlike.Client) {
	var (
		err error
		ctx = context.Background()
	)

	{
		dbName := "test"
		err = client.DropDatabase(ctx, dbName)
		require.NoError(t, err)
		err = client.CreateDatabase(ctx, dbName)
		require.NoError(t, err)

		testDB := client.Database(dbName)
		require.Equal(t, dbName, testDB.Name())
		tb := testDB.Table("t1")
		require.NotNil(t, tb)

		tb.MustUnsafeMigrate(ctx, struct {
			ID int64 `sqlike:",auto_increment"`
		}{})

		var o struct {
			Rows uint `sqlike:"rows"`
		}

		err = testDB.QueryRow(ctx, "EXPLAIN SELECT * FROM `t1`;").Decode(&o)
		require.NoError(t, err)
		// empty table will treat as one record when using EXPLAIN
		require.Equal(t, uint(1), o.Rows)
	}

	{
		dbs, err := client.ListDatabases(ctx)
		require.True(t, len(dbs) > 0)
		require.NoError(t, err)
	}

}
