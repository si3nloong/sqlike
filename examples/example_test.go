package examples

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/stretchr/testify/require"
)

// TestExamples :
func TestExamples(t *testing.T) {
	client, err := sqlike.Connect("mysql",
		options.Connect().
			SetUsername(`root`).
			SetPassword(`abcd1234`).
			SetDatabase(`sqlike`))
	require.NoError(t, err)
	dbs, err := client.ListDatabases()
	require.NoError(t, err)
	log.Println(dbs)
	db := client.Database("sqlike")

	MigrateExamples(t, db)
	InsertExamples(t, db)
	FindExamples(t, db)
	UpdateExamples(t, db)
	DeleteExamples(t, db)
}
