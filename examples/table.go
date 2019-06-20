package examples

import (
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/stretchr/testify/require"
)

// MigrateExamples :
func MigrateExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns  *normalStruct
		err error
	)

	table := db.Table("normal_struct")
	{
		err = db.Table("NormalStruct").DropIfExits()
		require.NoError(t, err)
	}

	{
		err = table.Migrate(ns)
		require.NoError(t, err)
	}

	{
		err = table.Rename("NormalStruct")
		require.NoError(t, err)
	}

	{
		err = db.Table("NormalStruct").Truncate()
		require.NoError(t, err)
	}

	// Alter table
	{
		err = db.Table("NormalStruct").Migrate(normalStruct{})
		require.NoError(t, err)
	}
}

// MigrateErrorExamples :
func MigrateErrorExamples(t *testing.T, db *sqlike.Database) {
	var (
		err error
	)

	{
		// empty table shouldn't able to migrate
		err = db.Table("").Migrate(new(normalStruct))
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(int(1))
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(struct{}{})
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(nil)
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(bool(false))
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(map[string]interface{}{})
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate([]interface{}{})
		require.Error(t, err)
	}
}
