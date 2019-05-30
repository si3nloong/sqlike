package examples

import (
	"database/sql"
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

// MigrateExamples :
func MigrateExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns  *normalStruct
		err error
	)

	table := db.Table("NormalStruct")
	{
		err = table.Migrate(ns)
		require.NoError(t, err)
	}
	{
		err = table.Truncate()
		require.NoError(t, err)
	}

}

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
		result, err = table.InsertOne(&ns)
		require.NoError(t, err)
		affected, err = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	}

	// Single upsert
	{
		ns.Emoji = `🤕`
		result, err = table.InsertOne(&ns,
			options.InsertOne().
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
		result, err = table.InsertMany(&nss)
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
