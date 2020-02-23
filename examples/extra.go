package examples

import (
	"context"
	"testing"

	"github.com/si3nloong/sqlike/sql"
	"github.com/si3nloong/sqlike/sqlike"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

// ExtraExamples :
func ExtraExamples(t *testing.T, db *sqlike.Database, mg *mongo.Database) {
	var (
		err error
		ctx = context.Background()
	)

	table := db.Table("A")

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
