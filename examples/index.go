package examples

import (
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
	)

	table := db.Table("Index")

	{

		err = table.DropIfExits()
		require.NoError(t, err)
	}

	// Migrate and create unique index with `unique_index` tag
	{
		err = table.Migrate(indexStruct{})
		require.NoError(t, err)
	}

	// Create one index
	{
		idx := table.Indexes()
		err = idx.CreateOne(indexes.Index{
			Columns: indexes.Columns("ID"),
		})
		require.NoError(t, err)
		idxs, err = idx.List()
		require.True(t, len(idxs) > 1)
		require.NoError(t, err)
	}

	// Auto build indexes using yaml file
	{
		err = db.BuildIndexes()
		require.NoError(t, err)
		idxs, err = db.Table("NormalStruct").Indexes().List()
		require.NoError(t, err)
		require.Contains(t, idxs, sqlike.Index{
			// Name:      "IX-SID@ASC;Emoji@ASC;Bool@DESC",
			Name:      "eb8bc4a93ee6af77e2ec575e12935e6d",
			Type:      "BTREE",
			IsVisible: true,
		})
		require.Contains(t, idxs, sqlike.Index{
			Name:      "test_idx",
			Type:      "BTREE",
			IsVisible: true,
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
		err = iv.Create(idxs)
		require.NoError(t, err)
		ok, _ = table.HasIndexByName("Bool_Int")
		require.True(t, ok)
		ok, _ = table.HasIndexByName("DateTime_Timestamp")
		require.True(t, ok)
		err = iv.CreateIfNotExists(idxs)
		require.NoError(t, err)
	}

}
