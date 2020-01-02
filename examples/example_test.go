package examples

import (
	"log"
	"testing"

	"github.com/Masterminds/semver"
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
}

// TestExamples :
func TestExamples(t *testing.T) {
	client, err := sqlike.Connect("mysql",
		options.Connect().
			SetUsername("root").
			SetPassword("abcd1234").
			SetCharset("utf8mb4"),
	)
	require.NoError(t, err)

	mg := connectMongoDB()

	v := client.Version()
	require.Equal(t, "mysql", client.DriverName())
	require.True(t, v.GreaterThan(semver.MustParse("5.7")))

	dbs, err := client.ListDatabases()
	require.True(t, len(dbs) > 0)
	require.NoError(t, err)

	db := client.SetLogger(Logger{}).
		Database("sqlike")

	{
		MigrateExamples(t, db)
		IndexExamples(t, db)
		InsertExamples(t, db)
		FindExamples(t, db)
		TransactionExamples(t, db)
		PaginationExamples(t, client)
		UpdateExamples(t, db)
		DeleteExamples(t, db)
		ExtraExamples(t, db, mg)
		JSONExamples(t, db)
		CasbinExamples(t, db)
		SpatialExamples(t, db)
	}

	// Errors
	{
		MigrateErrorExamples(t, db)
		InsertErrorExamples(t, db)
		FindErrorExamples(t, db)
	}

}
