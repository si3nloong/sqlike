package examples

import (
	"context"
	"testing"

	"github.com/si3nloong/sqlike"
	"github.com/stretchr/testify/require"
)

// MigrateExamples :
func MigrateExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
	var (
		ns      *normalStruct
		err     error
		results []sqlike.Column
		columns []string
	)

	table := db.Table("normal_struct")

	{
		err = db.Table("NormalStruct").DropIfExists(ctx)
		require.NoError(t, err)
	}

	// migrate table
	{
		err = table.UnsafeMigrate(ctx, ns)
		require.NoError(t, err)
		columnMap := make(map[string]sqlike.Column)
		columns = make([]string, 0)
		results, err = table.Columns().List(ctx)
		require.NoError(t, err)

		table.MustMigrate(ctx, ns)

		for _, f := range results {
			columnMap[f.Name] = f
			columns = append(columns, f.Name)
		}

		{
			pk := columnMap["$Key"]
			require.Equal(t, "VARCHAR(36)", pk.Type)
			require.Equal(t, "$Key", pk.Name)
			require.Equal(t, "Primary key", pk.Comment)
		}

		// check struct tag option
		require.Equal(t, "VARCHAR(300)", columnMap["CustomStrType"].Type)
		require.Equal(t, "DOUBLE UNSIGNED", columnMap["UFloat32"].Type)
		require.Equal(t, "ENUM('SUCCESS','FAILED','UNKNOWN')", columnMap["Enum"].Type)
		// // enum by default is using latin1 for performance concern
		// require.Equal(t, &latin1, columnMap["Enum"].Charset)

		require.ElementsMatch(t, []string{
			"$Key", "Key", "PtrUUID", "Date", "SID",
			"Emoji", "FullText", "LongStr", "CustomStrType",
			"EmptyByte", "Byte", "Bool",
			"Int", "TinyInt", "SmallInt", "MediumInt", "BigInt",
			"Uint", "TinyUint", "SmallUint", "MediumUint", "BigUint",
			"Float32", "Float64", "UFloat32",
			"EmptyStruct",
			"Struct", "VirtualColumn", "Struct.StoredStr",
			"JSONRaw", "Map",
			"DateTime", "Timestamp",
			"Language", "Languages", "Currency", "Currencies",
			"Enum", "Set", "CreatedAt", "UpdatedAt",
		}, columns)
	}

	{
		err = table.Rename(ctx, "NormalStruct")
		require.NoError(t, err)
	}

	{
		err = db.Table("NormalStruct").Truncate(ctx)
		require.NoError(t, err)
	}

	// Alter table
	{
		err = db.Table("NormalStruct").Migrate(ctx, normalStruct{})
		require.NoError(t, err)
	}

	{
		err = db.Table("PtrStruct").DropIfExists(ctx)
		require.NoError(t, err)
	}

	{
		err = db.Table("PtrStruct").Migrate(ctx, ptrStruct{})
		require.NoError(t, err)
	}

	// migrate table with generated columns
	{
		table := db.Table("GeneratedStruct")
		err = table.DropIfExists(ctx)
		require.NoError(t, err)

		err = table.Migrate(ctx, generatedStruct{})
		require.NoError(t, err)

		columns = make([]string, 0)
		results, err = table.Columns().List(ctx)
		require.NoError(t, err)

		for _, f := range results {
			columns = append(columns, f.Name)
		}

		require.ElementsMatch(t, []string{
			"NestedID", "Amount", "Nested", "CivilDate",
			"No", "id", "Line1", "Line2", "City", "State", "Country",
			"Date.CreatedAt", "Date.UpdatedAt",
		}, columns)

		err = table.Migrate(ctx, generatedStruct{})
		require.NoError(t, err)
	}

	temp := db.Table("Temp")

	// migrate with auto_increment field
	{
		err = temp.DropIfExists(ctx)
		require.NoError(t, err)
		temp.MustMigrate(ctx,
			struct {
				ID     string `sqlike:"$Key"`
				Number int64  `sqlike:",auto_increment"`
			}{},
		)
	}

	// migrate with auto_increment field
	{
		err = temp.DropIfExists(ctx)
		require.NoError(t, err)
		temp.MustMigrate(ctx, struct {
			ID     string `sqlike:"$Key"`
			Number int64
		}{})

		temp.MustMigrate(ctx, struct {
			ID     string `sqlike:"$Key"`
			Number int64  `sqlike:",auto_increment"`
		}{})
	}
}

// MigrateErrorExamples :
func MigrateErrorExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {
	var (
		err error
	)

	{
		// empty table shouldn't able to migrate
		err = db.Table("").Migrate(ctx, new(normalStruct))
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(ctx, int(1))
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(ctx, struct{}{})
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(ctx, nil)
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(ctx, bool(false))
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(ctx, map[string]interface{}{})
		require.Error(t, err)

		err = db.Table("NormalStruct").Migrate(ctx, []interface{}{})
		require.Error(t, err)
	}
}
