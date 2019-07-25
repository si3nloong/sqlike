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
				SetOmitFields("Int").
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
		m := make(map[string]int)
		m["one"] = 1
		m["two"] = 2
		ns.Map = m
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

// InsertErrorExamples :
func InsertErrorExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns  normalStruct
		err error
	)

	{
		_, err = db.Table("NormalStruct").InsertOne(nil)
		require.Error(t, err)

		var uninitialized *normalStruct
		_, err = db.Table("NormalStruct").InsertOne(uninitialized)
		require.Error(t, err)

		ns = normalStruct{}
		_, err = db.Table("NormalStruct").InsertOne(ns)
		require.Error(t, err)

		ns = normalStruct{}
		_, err = db.Table("NormalStruct").InsertOne(&ns)
		require.Error(t, err)
	}

	{
		_, err = db.Table("NormalStruct").InsertMany(
			[]normalStruct{},
		)
		require.Error(t, err)
	}
}
