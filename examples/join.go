package examples

import (
	"context"
	"testing"

	"github.com/si3nloong/sqlike/v2"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/stretchr/testify/require"
)

// JoinExamples :
func JoinExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {

	t.Run("LeftJoin", func(t *testing.T) {
		result, err := db.QueryStmt(ctx, sql.Select().
			From(db.Name(), "User").
			LeftJoin(),
		)
		require.NoError(t, err)
		defer result.Close()
	})

	t.Run("InnerJoin", func(t *testing.T) {
		result, err := db.QueryStmt(ctx, sql.Select().
			From(db.Name(), "User").
			LeftJoin(),
		)
		require.NoError(t, err)
		defer result.Close()
	})

	// panic("")

}
