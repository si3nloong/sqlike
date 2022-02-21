package examples

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/si3nloong/sqlike/v2"
	"github.com/si3nloong/sqlike/v2/actions"
	"github.com/si3nloong/sqlike/v2/options"
	sqldump "github.com/si3nloong/sqlike/v2/sql/dump"
	"github.com/si3nloong/sqlike/v2/sql/expr"
	"github.com/si3nloong/sqlike/v2/types"
	"github.com/stretchr/testify/require"
)

type dumpStruct struct {
	UUID    uuid.UUID
	String  string
	Bool    bool
	Int64   int64
	Int     int
	Uint64  uint64
	Uint    uint
	Byte    []byte
	JSONRaw json.RawMessage
	JSON    struct{}
	Array   []string
	// Point        orb.Point
	// LineString   orb.LineString
	Enum        Enum      `sqlike:",enum=SUCCESS|FAILED|UNKNOWN"`
	Set         types.Set `sqlike:",set=A|B|C"`
	Date        civil.Date
	DateTime    time.Time
	PtrString   *string
	PtrBool     *bool
	PtrInt64    *int64
	PtrUint64   *uint64
	PtrJSONRaw  *json.RawMessage
	PtrDate     *civil.Date
	PtrTime     *civil.Time
	PtrDateTime *time.Time
	PtrSet      *types.Set `sqlike:",set=A|B|C"`
}

// SQLDumpExamples :
func SQLDumpExamples(ctx context.Context, t *testing.T, client *sqlike.Client) {

	db := client.Database("sqlike")
	table := db.Table("sqldump")

	if err := table.DropIfExists(ctx); err != nil {
		require.NoError(t, err)
	}
	table.MustUnsafeMigrate(ctx, dumpStruct{})

	// generate data
	data := make([]dumpStruct, 10)
	for i := 0; i < len(data); i++ {
		data[i] = newDumpStruct()
	}

	{
		if _, err := table.Insert(
			ctx,
			&data,
			options.Insert().SetDebug(true),
		); err != nil {
			require.NoError(t, err)
		}
	}

	{
		// zip all sql files

		file, err := ioutil.TempFile("", ".sql")
		if err != nil {
			panic(err)
		}
		defer os.Remove(file.Name())

		// utcNow := time.Now().UTC().Add(-1 * time.Hour * 48)
		zero := time.Time{}
		filter := expr.And(
			expr.GreaterThan("DateTime", zero.Format("2006-01-02 15:04:05")),
		)

		offset := uint(0)
		limit := uint(10)
		dumper := sqldump.NewDumper("mysql", client)

		// backup 100 records per time
		for {
			// check how many records return or backup
			affected, err := dumper.BackupTo(
				ctx,
				actions.Find().
					From("sqlike", "sqldump").
					Where(filter).
					Offset(offset).
					Limit(limit),
				file,
			)
			if err != nil {
				if err != sql.ErrNoRows {
					require.NoError(t, err)
				}
				break
			}

			if affected < int64(limit) {
				break
			}

			offset += limit
		}

		file.Close()

		err = table.Truncate(ctx)
		require.NoError(t, err)

		// b, err := ioutil.ReadFile(file.Name())
		// if err != nil {
		// 	panic(err)
		// }

		// log.Println(string(b))
		// result, err := client.Exec(string(b))
		// require.NoError(t, err)
		// affected, err := result.RowsAffected()
		// require.NoError(t, err)
		// log.Println(affected)
	}
}

func newDumpStruct() (o dumpStruct) {
	date := gofakeit.Date()
	o.UUID = uuid.New()
	o.String = gofakeit.Name()
	o.Int = int(gofakeit.Int32())
	o.Int64 = gofakeit.Int64()
	o.Uint = uint(gofakeit.Uint32())
	o.Uint64 = gofakeit.Uint64()
	o.JSONRaw = []byte(`{"id": 100, "message": "hello world"}`)
	o.Date = civil.Date{Year: date.Year(), Month: date.Month(), Day: date.Day()}
	o.Enum = Enum(gofakeit.RandString([]string{
		"SUCCESS",
		"FAILED",
		"UNKNOWN",
	}))
	o.Set = []string{"A", "C"}
	o.DateTime = gofakeit.Date()
	// o.Timestamp = types.Timestamp(gofakeit.Date())
	// o.Point = orb.Point{gofakeit.Longitude(), gofakeit.Latitude()}
	// o.LineString = orb.LineString{
	// 	orb.Point{gofakeit.Longitude(), gofakeit.Latitude()},
	// 	orb.Point{gofakeit.Longitude(), gofakeit.Latitude()},
	// 	orb.Point{gofakeit.Longitude(), gofakeit.Latitude()},
	// }

	str := gofakeit.Address().Address
	num := gofakeit.Int64()
	ts := gofakeit.Date()

	o.PtrString = &str
	o.PtrInt64 = &num
	o.PtrDateTime = &ts
	return
}
