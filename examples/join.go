package examples

import (
	"context"
	"testing"

	"github.com/si3nloong/sqlike/v2"
	"github.com/si3nloong/sqlike/v2/sql"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/stretchr/testify/require"
)

// JoinExamples :
func JoinExamples(ctx context.Context, t *testing.T, db *sqlike.Database) {

	// setup databases
	db.Table("User").MustMigrate(ctx, User{})
	addressTable := db.Table("UserAddress")
	addressTable.MustMigrate(ctx, UserAddress{})

	t.Run("Check the foreign key index", func(t *testing.T) {
		idxs, err := addressTable.Indexes().List(ctx)
		require.NoError(t, err)

		require.Contains(t, idxs, sqlike.Index{Name: "UserID", Type: "BTREE"})
	})

	// SELECT * FROM `User` LEFT JOIN `UserAddress` ON `User`.`ID` = `UserAddress`.`UserID`;
	t.Run("LeftJoin", func(t *testing.T) {
		result, err := db.QueryStmt(ctx, sql.
			Select(
				expr.Pair("u", "ID"),
				expr.Pair("u", "Name"),
				expr.Pair("u", "Age"),
				expr.Pair("u", "Status"),
				expr.Pair("u", "CreatedAt"),
				expr.As(expr.Pair("ua", "ID"), "UserAddressID"),
				expr.Pair("ua", "UserID"),
			).
			From(expr.As(expr.Pair("sqlike", "User"), "u")).
			LeftJoin(
				expr.As(expr.Pair("sqlike", "UserAddress"), "ua"),
				expr.Pair("u", "ID"),
				expr.Pair("ua", "UserID"),
			),
		)
		require.NoError(t, err)
		defer result.Close()

		require.ElementsMatch(t, []string{"ID", "Name", "Age", "Status", "CreatedAt", "UserAddressID", "UserID"}, result.Columns())

		for result.Next() {

		}
	})

	// SELECT * FROM `User` INNER JOIN `UserAddress` ON `User`.`ID` = `UserAddress`.`UserID`;
	t.Run("InnerJoin", func(t *testing.T) {
		result, err := db.QueryStmt(ctx, sql.
			Select(
				expr.Pair("u", "ID"),
				expr.Pair("u", "Name"),
				expr.Pair("u", "Age"),
				expr.Pair("u", "Status"),
				expr.Pair("u", "CreatedAt"),
				expr.As(expr.Pair("ua", "ID"), "UserAddressID"),
				expr.Pair("ua", "UserID"),
			).
			From(expr.As(expr.Pair("sqlike", "User"), "u")).
			InnerJoin(
				expr.As(expr.Pair("sqlike", "UserAddress"), "ua"),
				expr.Pair("u", "ID"),
				expr.Pair("ua", "UserID"),
			),
		)
		require.NoError(t, err)
		defer result.Close()

		for result.Next() {

		}
	})

}
