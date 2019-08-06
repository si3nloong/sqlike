package examples

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

type Logger struct {
}

func (l Logger) Debug(stmt *sqlstmt.Statement) {
	// log.Printf("%v", stmt)
	log.Printf("%+v", stmt)
	return
}

// TestExamples :
func TestExamples(t *testing.T) {
	client, err := sqlike.Connect("mysql",
		options.Connect().
			SetUsername(`root`).
			SetPassword(`abcd1234`),
	)
	require.NoError(t, err)

	db := client.SetLogger(Logger{}).
		Database("sqlike")

	dbs, err := client.ListDatabases()
	require.True(t, len(dbs) > 0)
	require.NoError(t, err)
	// client.SetLogger()
	log.Println(dbs)

	// Errors
	{
		MigrateErrorExamples(t, db)
		InsertErrorExamples(t, db)
		FindErrorExamples(t, db)
	}

	MigrateExamples(t, db)
	IndexExamples(t, db)
	InsertExamples(t, db)
	FindExamples(t, db)
	TransactionExamples(t, db)
	PaginationExamples(t, client)
	UpdateExamples(t, db)
	DeleteExamples(t, db)
	ExtraExamples(t, db)

}
