package examples

import (
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
	"github.com/stretchr/testify/require"
)

// UpdateExamples :
func UpdateExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns       normalStruct
		err      error
		affected int64
	)

	table := db.Table("NormalStruct")

	{
		err = table.FindOne(nil).Decode(&ns)
		require.NoError(t, err)
		affected, err = table.UpdateOne(
			actions.UpdateOne().
				Where(expr.Equal("$Key", ns.ID)).
				Set("Emoji", "<😗>"))
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	}

}
