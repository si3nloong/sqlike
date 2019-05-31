package examples

import (
	"log"
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
	"github.com/stretchr/testify/require"
)

// FindExamples :
func FindExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns  normalStruct
		err error
	)

	table := db.Table("NormalStruct")
	{
		err = table.FindOne(
			actions.FindOne().Where(
				expr.Equal("$Key", "1000"),
			),
		).Decode(&ns)
		log.Println(err)
		require.Error(t, err)
	}
}
