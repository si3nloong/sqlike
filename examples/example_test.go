package examples

import (
	"context"
	"log"
	"testing"

	"github.com/Masterminds/semver/v3"

	mysql "github.com/go-sql-driver/mysql"

	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/plugin/opentracing"
	"github.com/si3nloong/sqlike/v2/sql/instrumented"

	"github.com/si3nloong/sqlike/v2"
	"github.com/stretchr/testify/require"
)

type Logger struct {
}

func (l Logger) Debug(args ...any) {
	// log.Printf("%v", stmt)
	log.Printf("%+v", args...)
}

// TestExamples :
func TestExamples(t *testing.T) {
	var (
		ctx = context.Background()
	)

	// normal connect
	{
		client := sqlike.MustConnect(
			ctx,
			"mysql",
			options.Connect().
				ApplyURI(`root:abcd1234@tcp()/sqlike?parseTime=true&multiStatements=true&loc=UTC&charset=utf8mb4&collation=utf8mb4_general_ci`),
		)

		// set timezone for UTC
		if _, err := client.ExecContext(ctx, `SET GLOBAL time_zone = '+00:00';`); err != nil {
			panic(err)
		}

		testCase(ctx, t, client)
	}

	// with tracing (OpenTracing)
	{
		driver := "mysql"
		username := "root"

		cfg := mysql.NewConfig()
		cfg.User = username
		cfg.Params = map[string]string{
			"charset": "utf8mb4",
			// "collation": "utf8mb4_general_ci",
		}
		cfg.Collation = "utf8mb4_general_ci"
		cfg.Passwd = "abcd1234"
		cfg.ParseTime = true
		conn, err := mysql.NewConnector(cfg)
		if err != nil {
			panic(err)
		}

		itpr := opentracing.NewInterceptor(
			opentracing.WithDBInstance("sqlike"),
			opentracing.WithDBUser(username),
			opentracing.WithDBType(driver),
			opentracing.WithExec(true),
			opentracing.WithQuery(true),
		)
		client := sqlike.MustConnectDB(ctx, driver, instrumented.WrapConnector(conn, itpr))
		defer client.Close()
		testCase(ctx, t, client)
	}

}

func testCase(ctx context.Context, t *testing.T, client *sqlike.Client) {
	v := client.Version()
	require.Equal(t, "mysql", client.DriverName())
	require.True(t, v.GreaterThan(semver.MustParse("5.7")))
	client.SetLogger(Logger{})
	DatabaseExamples(t, client)
	db := client.Database("sqlike")
	mg := connectMongoDB(ctx)

	{
		MigrateExamples(ctx, t, db)
		IndexExamples(ctx, t, db)

		InsertExamples(ctx, t, db)
		FindExamples(ctx, t, db)
		JoinExamples(ctx, t, db)
		AdvanceQueryExamples(ctx, t, db)
		TransactionExamples(ctx, t, db)
		PaginationExamples(ctx, t, client)
		UpdateExamples(ctx, t, db)
		DeleteExamples(ctx, t, db)
		JSONExamples(ctx, t, db)
		CasbinExamples(ctx, t, db)
		// SpatialExamples(ctx, t, db)

		// SQLDumpExamples(ctx, t, client)

		ExtraExamples(ctx, t, db, mg)
	}

	// Errors
	{
		MigrateErrorExamples(ctx, t, db)
		InsertErrorExamples(ctx, t, db)
		FindErrorExamples(ctx, t, db)
		UpdateErrorExamples(ctx, t, db)
	}
}
