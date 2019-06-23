package examples

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/options"
	sqlstmt "github.com/si3nloong/sqlike/sqlike/sql/stmt"
	"github.com/stretchr/testify/require"
)

type Logger struct {
}

func (l Logger) Format(stmt *sqlstmt.Statement) {
	log.Println("Debug here ======================>")
	log.Printf("%v", stmt)
	log.Printf("%+v", stmt)
	return
}

// TestExamples :
func TestExamples(t *testing.T) {
	client, err := sqlike.Connect("mysql",
		options.Connect().
			SetUsername(`root`).
			SetPassword(`abcd1234`).
			SetDatabase(`sqlike`),
	)
	require.NoError(t, err)

	dbs, err := client.ListDatabases()
	require.NoError(t, err)
	// client.SetLogger()
	log.Println(dbs)

	db := client.SetLogger(Logger{}).
		Database("sqlike")

	MigrateExamples(t, db)
	InsertExamples(t, db)
	FindExamples(t, db)
	TransactionExamples(t, db)
	return
	PaginationExamples(t, db)
	UpdateExamples(t, db)
	DeleteExamples(t, db)

	MigrateErrorExamples(t, db)
	InsertErrorExamples(t, db)
	FindErrorExamples(t, db)

}
