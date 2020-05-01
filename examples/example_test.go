package examples

import (
	"context"
	"log"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/go-sql-driver/mysql"

	"github.com/si3nloong/sqlike/plugin/opentracing"
	"github.com/si3nloong/sqlike/sql/instrumented"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"

	"github.com/si3nloong/sqlike/sqlike"
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
	var (
		ctx = context.Background()
	)

	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = "abcd1234"
	cfg.ParseTime = true
	conn, err := mysql.NewConnector(cfg)
	if err != nil {
		panic(err)
	}

	itpr := opentracing.Interceptor(
		opentracing.WithDBInstance("sqlike"),
		opentracing.WithDBUser("root"),
		opentracing.WithExec(true),
		opentracing.WithQuery(true),
	)
	client, err := sqlike.ConnectDB(ctx, "mysql", instrumented.WrapConnector(conn, itpr))
	require.NoError(t, err)
	defer client.Close()

	// client, err := sqlike.Connect(
	// 	ctx,
	// 	"mysql",
	// 	options.Connect().
	// 		SetUsername("root").
	// 		SetPassword("abcd1234").
	// 		SetCharset("utf8mb4"),
	// )

	mg := connectMongoDB()

	v := client.Version()
	require.Equal(t, "mysql", client.DriverName())
	require.True(t, v.GreaterThan(semver.MustParse("5.7")))

	// client.SetLogger(Logger{})
	DatabaseExamples(t, client)
	db := client.Database("sqlike")

	{
		MigrateExamples(t, ctx, db)
		IndexExamples(t, ctx, db)

		InsertExamples(t, ctx, db)
		FindExamples(t, ctx, db)
		QueryExamples(t, ctx, db)
		TransactionExamples(t, ctx, db)
		PaginationExamples(t, ctx, client)
		UpdateExamples(t, ctx, db)
		DeleteExamples(t, ctx, db)
		JSONExamples(t, ctx, db)
		CasbinExamples(t, ctx, db)
		SpatialExamples(t, ctx, db)
		ExtraExamples(t, ctx, db, mg)
	}

	// Errors
	{
		MigrateErrorExamples(t, db)
		InsertErrorExamples(t, db)
		FindErrorExamples(t, db)
		UpdateErrorExamples(t, db)
	}

}
