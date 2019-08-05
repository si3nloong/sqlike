package examples

import (
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/stretchr/testify/require"
)

// ExtraExamples :
func ExtraExamples(t *testing.T, db *sqlike.Database) {
	var (
		err error
	)

	table := db.Table("A")

	{
		table.MustMigrate(normalStruct{})
		table.Truncate()

		err = table.Copy([]string{
			"$Key", "SID", "Emoji", "LongStr",
			"TinyInt", "Float64", "EmptyStruct", "Struct",
		}, actions.Copy().
			From("sqlike", "NormalStruct").
			Select(
				"$Key", "SID", "Emoji", "LongStr",
				"TinyInt", "Float32", "EmptyStruct", "Struct",
			),
		)
		require.NoError(t, err)
	}
}
