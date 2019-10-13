package examples

import (
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/stretchr/testify/require"
)

// MigrateExamples :
func MigrateExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns      *normalStruct
		err     error
		results []sqlike.Column
		columns []string
	)

	{
		err = db.DropIfExists()
		require.NoError(t, err)

		err = db.Drop()
		require.Error(t, err)

		err = db.CreateIfNotExists()
		require.NoError(t, err)

		err = db.Create()
		require.Error(t, err)
	}

	table := db.Table("normal_struct")

	{
		err = db.Table("NormalStruct").DropIfExits()
		require.NoError(t, err)
	}

	// migrate NormalStruct
	{
		err = table.Migrate(ns)
		require.NoError(t, err)

		columnMap := make(map[string]sqlike.Column)
		columns = make([]string, 0)
		results, err = table.Columns().List()
		require.NoError(t, err)

		for _, f := range results {
			columnMap[f.Name] = f
			columns = append(columns, f.Name)
		}

		// check struct tag option
		require.Equal(t, "VARCHAR(300)", columnMap["CustomStrType"].Type)
		require.Equal(t, "DOUBLE UNSIGNED", columnMap["UFloat32"].Type)
		require.Equal(t, "ENUM('SUCCESS','FAILED','UNKNOWN')", columnMap["Enum"].Type)

		require.ElementsMatch(t, []string{
			"$Key", "Key", "Date", "SID",
			"Emoji", "FullText", "LongStr", "CustomStrType",
			"EmptyByte", "Byte", "Bool",
			"Int", "TinyInt", "SmallInt", "MediumInt", "BigInt",
			"Uint", "TinyUint", "SmallUint", "MediumUint", "BigUint",
			"Float32", "Float64", "UFloat32",
			"EmptyStruct", "GeoPoint",
			"Struct", "VirtualColumn", "Struct.StoredStr",
			"JSONRaw", "Map",
			"DateTime", "Timestamp",
			"Language", "Languages", "Currency", "Currencies",
			"Enum", "CreatedAt", "UpdatedAt",
		}, columns)
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

	{
		err = db.Table("PtrStruct").DropIfExits()
		require.NoError(t, err)
	}

	{
		err = db.Table("PtrStruct").Migrate(ptrStruct{})
		require.NoError(t, err)
	}

	{
		table := db.Table("GeneratedStruct")
		err = table.DropIfExits()
		require.NoError(t, err)

		err = table.Migrate(generatedStruct{})
		require.NoError(t, err)

		columns = make([]string, 0)
		results, err = table.Columns().List()
		require.NoError(t, err)

		for _, f := range results {
			columns = append(columns, f.Name)
		}

		require.ElementsMatch(t, []string{
			"NestedID", "Amount", "Nested",
			"No", "id",
			"Line1", "Line2", "City", "State", "Country",
			"Date.CreatedAt", "Date.UpdatedAt",
		}, columns)

		err = table.Migrate(generatedStruct{})
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
