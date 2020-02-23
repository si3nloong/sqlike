package examples

import (
	"context"
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/stretchr/testify/require"
)

// IndexExamples :
func IndexExamples(t *testing.T, db *sqlike.Database) {
	var (
		err  error
		idxs []sqlike.Index
		ok   bool
		ctx  = context.Background()
	)

	table := db.Table("Index")

	{

		err = table.DropIfExists(ctx)
		require.NoError(t, err)
	}

	// Migrate and create unique index with `unique_index` tag
	{
		err = table.Migrate(ctx, indexStruct{})
		require.NoError(t, err)
	}

	// Create one index
	{
		idx := table.Indexes()
		err = idx.CreateOne(
			ctx,
			indexes.Index{
				Columns: indexes.Columns("ID"),
			})
		require.NoError(t, err)
		idxs, err = idx.List(ctx)
		require.True(t, len(idxs) > 1)
		require.NoError(t, err)
	}

	// Auto build indexes using yaml file
	{
		err = db.BuildIndexes(ctx)
		require.NoError(t, err)
		idxs, err = db.Table("NormalStruct").Indexes().List(ctx)
		require.NoError(t, err)
		require.Contains(t, idxs, sqlike.Index{
			// Name:      "IX-SID@ASC;Emoji@ASC;Bool@DESC",
			Name:     "eb8bc4a93ee6af77e2ec575e12935e6d",
			Type:     "BTREE",
			IsUnique: false,
		})
		require.Contains(t, idxs, sqlike.Index{
			Name:     "test_idx",
			Type:     "BTREE",
			IsUnique: false,
		})
	}

	table = db.Table("NormalStruct")

	// Create multiple indexes
	{
		idxs := []indexes.Index{
			indexes.Index{
				Name:    "Bool_Int",
				Type:    indexes.BTree,
				Columns: indexes.Columns("Bool", "Int"),
			},
			indexes.Index{
				Name:    "DateTime_Timestamp",
				Type:    indexes.BTree,
				Columns: indexes.Columns("DateTime", "Timestamp"),
			},
		}

		iv := table.Indexes()
		err = iv.Create(ctx, idxs)
		require.NoError(t, err)
		ok, _ = table.HasIndexByName(ctx, "Bool_Int")
		require.True(t, ok)
		ok, _ = table.HasIndexByName(ctx, "DateTime_Timestamp")
		require.True(t, ok)
		err = iv.CreateIfNotExists(ctx, idxs)
		require.NoError(t, err)
	}

}
