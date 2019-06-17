package examples

import (
	"database/sql"
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

// InsertExamples :
func InsertExamples(t *testing.T, db *sqlike.Database) {
	var (
		err      error
		result   sql.Result
		affected int64
	)

	table := db.Table("NormalStruct")
	ns := newNormalStruct()

	// Single insert
	{
		result, err = table.InsertOne(&ns,
			options.InsertOne().
				SetDebug(true))
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

	}

	// Single upsert
	// - https://dev.mysql.com/doc/refman/8.0/en/insert-on-duplicate.html
	{
		ns.Emoji = `🤕`
		result, err = table.InsertOne(&ns,
			options.InsertOne().
				SetDebug(true).
				SetMode(options.InsertOnDuplicate))
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(2), affected)
	}

	// Multiple insert
	{
		nss := [...]normalStruct{
			newNormalStruct(),
			newNormalStruct(),
			newNormalStruct(),
		}
		result, err = table.InsertMany(&nss,
			options.InsertMany().
				SetDebug(true))
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(3), affected)
	}

	// Error insertion
	{
		_, err = table.InsertOne(&struct {
			Interface interface{}
		}{})
		require.Error(t, err)
		_, err = table.InsertOne(struct{}{})
		require.Error(t, err)
		var empty *struct{}
		_, err = table.InsertOne(empty)
		require.Error(t, err)

		_, err = table.InsertMany([]interface{}{})
		require.Error(t, err)
	}
}
