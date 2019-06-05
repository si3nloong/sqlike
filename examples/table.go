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
