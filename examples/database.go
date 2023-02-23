package examples

import (
	"context"
	"testing"

	"github.com/si3nloong/sqlike/v2"
	"github.com/stretchr/testify/require"
)

// DatabaseExamples :
func DatabaseExamples(ctx context.Context, t *testing.T, client *sqlike.Client) {
	var (
		err    error
		dbName = "test"
	)

	t.Run("DropDatabase", func(t *testing.T) {
		err = client.DropDatabase(ctx, dbName)
		require.NoError(t, err)
	})

	t.Run("CreateDatabase", func(t *testing.T) {
		err = client.CreateDatabase(ctx, dbName)
		require.NoError(t, err)
	})

	t.Run("Database", func(t *testing.T) {
		testDB := client.Database(dbName)
		require.Equal(t, dbName, testDB.Name())
		tb := testDB.Table("t1")
		require.NotNil(t, tb)

		tb.UnsafeMigrate(ctx, struct {
			ID int64 `sqlike:",auto_increment"`
		}{})

		// FIXME:
		// var o struct {
		// 	Rows uint `sqlike:"rows"`
		// }

		// err = testDB.QueryRow(ctx, "EXPLAIN SELECT * FROM `t1`;").Decode(&o)
		// require.NoError(t, err)

		// // panic("")
		// // empty table will treat as one record when using EXPLAIN
		// require.Equal(t, uint(1), o.Rows)
	})

	t.Run("ListDatabases", func(t *testing.T) {
		dbs, err := client.ListDatabases(ctx)
		require.True(t, len(dbs) > 0)
		require.NoError(t, err)
	})
}
