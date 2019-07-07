package examples

import (
	"log"
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/stretchr/testify/require"
)

// IndexExamples :
func IndexExamples(t *testing.T, db *sqlike.Database) {
	var (
		err error
	)

	table := db.Table("Index")
	idx := table.Indexes()

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
		err = idx.CreateOne(indexes.Index{
			Columns: []string{"ID"},
		})
		require.NoError(t, err)
	}

	{
		var idxs []sqlike.Index
		idxs, err = idx.List()
		require.NoError(t, err)
		log.Println(idxs)
	}
}
