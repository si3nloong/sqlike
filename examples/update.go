package examples

import (
	"database/sql"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/sqlike/sql/expr"
	"github.com/stretchr/testify/require"
)

// UpdateExamples :
func UpdateExamples(t *testing.T, db *sqlike.Database) {
	var (
		ns       normalStruct
		err      error
		result   sql.Result
		affected int64
	)

	table := db.Table("NormalStruct")

	{

		uid, _ := uuid.FromString(`be72fc34-917b-11e9-af91-6c96cfd87b17`)

		ns = normalStruct{}
		ns.ID = uid
		ns.Timestamp = time.Now()
		result, err = table.InsertOne(&ns)
		affected, _ = result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		affected, err = table.UpdateOne(
			actions.UpdateOne().
				Where(expr.Equal("$Key", uid)).
				Set(
					expr.Field("Emoji", "<😗>"),
				),
			options.UpdateOne().SetDebug(true),
		)

		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	}

}
