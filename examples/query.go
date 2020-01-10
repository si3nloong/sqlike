package examples

import (
	"log"
	"testing"

	"github.com/si3nloong/sqlike/sql"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/stretchr/testify/require"
)

func QueryExamples(t *testing.T, db *sqlike.Database) {
	// table := db.Table("query")

	stmt := expr.Union(
		sql.Select().
			From("sqlike", "GeneratedStruct").
			Where(
				expr.Equal("State", ""),
				expr.GreaterThan("Date.CreatedAt", "2006-01-02 16:00:00"),
				expr.GreaterOrEqual("No", 0),
				expr.Or(
					expr.Equal("State", ""),
					expr.NotEqual("State", ""),
					expr.GreaterOrEqual("State", ""),
				),
			).
			OrderBy(
				expr.Desc("NestedID"),
				expr.Asc("Amount"),
			).
			Limit(1).
			Offset(1),
		sql.Select().
			From("sqlike", "GeneratedStruct").
			Limit(1),
	)

	{
		// table := db.Table("GeneratedStruct")
		// err = table.Truncate()
		// require.NoError(t, err)

		result, err := db.QueryStmt(stmt)
		require.NoError(t, err)
		defer result.Close()

		gss := make([]*generatedStruct, 0)
		for result.Next() {
			gs := new(generatedStruct)
			if err := result.Decode(gs); err != nil {
				panic(err)
			}
			gss = append(gss, gs)
		}

		// TODO: add test
		log.Println(gss)
	}

	{
		if err := db.RunInTransaction(func(ctx sqlike.SessionContext) error {
			result, err := ctx.QueryStmt(stmt)
			if err != nil {
				return err
			}
			defer result.Close()
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
