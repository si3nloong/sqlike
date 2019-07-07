package examples

import (
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/stretchr/testify/require"
)

// DeleteExamples :
func DeleteExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns  normalStruct
		err error
	)

	table := db.Table("NormalStruct")

	{
		err = table.FindOne(
			actions.FindOne().OrderBy(expr.Desc("$Key")),
		).Decode(&ns)
		require.NoError(t, err)
		err = table.DestroyOne(&ns)
		require.NoError(t, err)
	}
}
