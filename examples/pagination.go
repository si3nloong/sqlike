package examples

import (
	"testing"

	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
	"github.com/stretchr/testify/require"
)

// PaginationExamples :
func PaginationExamples(t *testing.T, db *sqlike.Database) {
	var (
		nss    []normalStruct
		cursor *sqlike.Cursor
		err    error
	)

	limits := uint(2)

	cursor, err = db.Table("NormalStruct").
		Find(actions.Find().
			OrderBy(
				expr.Asc("$Key"),
			).
			Limit(limits + 1))
	require.NoError(t, err)

	err = cursor.All(&nss)
	require.NoError(t, err)
	length := len(nss)
	if uint(length) > limits {
		// csr := nss[length-1].ID.String()
		nss = nss[:length-1]
	}
}
