package examples

import (
	"context"
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
	client, err := sqlike.Connect(
		context.Background(),
		"mysql",
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

	client.SetLogger(Logger{})
	DatabaseExamples(t, client)
	db := client.Database("sqlike")

	{
		MigrateExamples(t, db)
		IndexExamples(t, db)
		InsertExamples(t, db)
		FindExamples(t, db)
		QueryExamples(t, db)
		TransactionExamples(t, db)
		PaginationExamples(t, client)
		UpdateExamples(t, db)
		DeleteExamples(t, db)
		JSONExamples(t, db)
		CasbinExamples(t, db)
		SpatialExamples(t, db)
		ExtraExamples(t, db, mg)
	}

	// Errors
	{
		MigrateErrorExamples(t, db)
		InsertErrorExamples(t, db)
		FindErrorExamples(t, db)
		UpdateErrorExamples(t, db)
	}

}
