package examples

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/si3nloong/sqlike/sql"
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
	"github.com/si3nloong/sqlike/types"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

// ExtraExamples :
func ExtraExamples(ctx context.Context, t *testing.T, db *sqlike.Database, mg *mongo.Database) {
	var (
		err error
	)

	table := db.Table("A")

	// replace into
	{
		table.MustMigrate(ctx, normalStruct{})
		err = table.Truncate(ctx)
		require.NoError(t, err)

		err = table.Replace(
			ctx,
			[]string{
				"$Key", "SID", "Date", "Emoji", "LongStr",
				"TinyInt", "Float64", "EmptyStruct", "Struct",
			},
			sql.Select(
				"$Key", "SID", "Date", "Emoji", "LongStr",
				"TinyInt", "Float32", "EmptyStruct", "Struct",
			).From("sqlike", "NormalStruct"),
		)
		require.NoError(t, err)
	}

	// set custom primary key
	{
		tbl := db.Table("PK")
		var a struct {
			Key string `sqlike:"$Key"`
			No  int64
		}

		err = tbl.DropIfExists(ctx)
		require.NoError(t, err)
		tbl.MustMigrate(ctx, a)

		var b struct {
			Key string `sqlike:"$Key"`
			No  int64  `sqlike:",primary_key"`
		}

		tbl.MustMigrate(ctx, b)
	}

	table = db.Table("B")

	// Alter table should add primary key if it's not exists
	{
		err = table.DropIfExists(ctx)
		require.NoError(t, err)
		table.MustMigrate(ctx, struct {
			ID   string
			Name string
		}{})

		table.MustMigrate(ctx, struct {
			ID   string `sqlike:",primary_key"`
			Name string
		}{})

		idxs, err := table.Indexes().List(ctx)
		require.NoError(t, err)
		require.Contains(t, idxs, sqlike.Index{
			Name:     "PRIMARY",
			Type:     "BTREE",
			IsUnique: true,
		})
	}

	// on update for datetime, time & date
	{
		table := db.Table("dateTime")
		if err := table.DropIfExists(ctx); err != nil {
			panic(err)
		}

		type dateTime struct {
			ID   uuid.UUID `sqlike:",primary_key"`
			Name string
			Date types.Date `sqlike:",on_update"`
			Time time.Time  `sqlike:",on_update"`
		}

		utcNow := time.Now().UTC()
		dt := dateTime{}
		dt.ID = uuid.New()
		dt.Name = "yuki"
		dt.Date = types.Date{}
		dt.Time = utcNow

		table.MustUnsafeMigrate(ctx, dt)

		if _, err := table.InsertOne(
			ctx,
			&dt,
			options.InsertOne().SetDebug(true),
		); err != nil {
			panic(err)
		}

		duration := time.Second * 5
		time.Sleep(duration)

		if _, err := table.UpdateOne(
			ctx,
			actions.UpdateOne().
				Where(
					expr.Equal("ID", dt.ID),
				).
				Set(
					expr.ColumnValue("Name", "eyoki"),
				),
		); err != nil {
			panic(err)
		}

		var o dateTime
		if err := table.FindOne(
			ctx,
			actions.FindOne().
				Where(
					expr.Equal("ID", dt.ID),
				),
		).Decode(&o); err != nil {
			panic(err)
		}

		require.True(t, o.Time.After(utcNow))
		require.Equal(t, o.Name, "eyoki")
		require.True(t, o.Time.Sub(utcNow) >= duration)

		type dateTimeNoUpdate struct {
			ID   uuid.UUID `sqlike:",primary_key"`
			Name string
			Date time.Time
			Time time.Time
		}
		table.MustUnsafeMigrate(ctx, dateTimeNoUpdate{})

		dtu := dateTimeNoUpdate{}
		dtu.ID = uuid.New()
		dtu.Name = "sianloong"
		dtu.Date = utcNow
		dtu.Time = utcNow
		if _, err := table.InsertOne(
			ctx,
			&dtu,
			options.InsertOne().SetDebug(true),
		); err != nil {
			panic(err)
		}
	}

	// MongoDB :
	// {
	// 	ctx := context.Background()
	// 	coll := mg.Collection("MongoStruct")
	// 	coll.Drop(ctx)

	// 	pk := types.NewNameKey("MongoStruct", types.NewIDKey("ID", nil))
	// 	msg := "hello world!!!"

	// 	ms := mongoStruct{}
	// 	ms.Key = pk
	// 	ms.Name = msg
	// 	_, err = coll.InsertOne(ctx, ms)

	// 	result := mongoStruct{}
	// 	err = coll.FindOne(ctx, bson.M{"key": pk}).
	// 		Decode(&result)
	// 	require.NoError(t, err)
	// 	require.Equal(t, pk, result.Key)
	// 	require.Equal(t, msg, result.Name)
	// }
}
