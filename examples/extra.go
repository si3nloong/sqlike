package examples

import (
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
	)

	table := db.Table("A")

	{
		table.MustMigrate(normalStruct{})
		err = table.Truncate()
		require.NoError(t, err)

		err = table.Replace(
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

		err = tbl.DropIfExists()
		require.NoError(t, err)
		tbl.MustMigrate(a)

		var b struct {
			Key string `sqlike:"$Key"`
			No  int64  `sqlike:",primary_key"`
		}

		tbl.MustMigrate(b)
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
