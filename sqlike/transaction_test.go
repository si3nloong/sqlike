package sqlike

import (
	"log"
	"math/rand"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/si3nloong/sqlike/sql/expr"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
)

type debug struct{}

// Debug :
func (l debug) Debug(stmt *sqlstmt.Statement) {
	log.Printf("%+v", stmt)
	return
}

func getMySQLClient() *Database {
	client := MustConnect("mysql",
		options.Connect().
			SetHost("127.0.0.1").
			SetPort("3306").
			SetUsername("root").
			SetPassword("test"),
	).
		SetPrimaryKey("ID").
		SetLogger(debug{})

	os.Stdout.WriteString("Connected to MySQL!")

	return client.Database("test")
}

func setupData(client *Database) *user {
	client.Table("user").Migrate(new(user))
	data := &user{ID: 1234, Name: "Oska"}
	client.Table("user").InsertOne(data)
	return data
}

func clearData(client *Database) {
	client.Table("user").DestroyOne(&user{ID: 1234})
}

type user struct {
	ID   int
	Name string
}

func TestCreateCommit(t *testing.T) {
	client := getMySQLClient()
	client.Table("user").Migrate(new(user))

	data := &user{ID: rand.Intn(10000), Name: "Oska"}
	trx, _ := client.BeginTransaction()
	trx.Table("user").InsertOne(data)

	if err := trx.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(&user{}); err != nil {
		t.Error("FindOne doesn't work well with InsertOne in transaction mode.")
	}

	if err := client.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(&user{}); err == nil {
		t.Error("InsertOne with transaciton mode shouldn't have result with normal FindOne.")
	}

	trx.CommitTransaction()

	if err := client.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(&user{}); err != nil {
		t.Error("Result should exists after rollback transaction.")
	}

	// Remove tested data
	client.Table("user").DestroyOne(data)
}

func TestCreateRollback(t *testing.T) {
	client := getMySQLClient()
	client.Table("user").Migrate(new(user))

	data := &user{ID: 1234, Name: "Oska"}
	trx, _ := client.BeginTransaction()
	trx.Table("user").InsertOne(data)

	if err := trx.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(&user{}); err != nil {
		t.Error("FindOne doesn't work well with InsertOne in transaction mode.")
	}

	if err := client.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(&user{}); err == nil {
		t.Error("InsertOne with transaciton mode shouldn't have result with normal FindOne.")
	}

	trx.RollbackTransaction()

	if err := client.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(&user{}); err == nil {
		t.Error("Result should exists after rollback transaction.")
	}
}

func TestUpdateCommit(t *testing.T) {
	client := getMySQLClient()
	data := setupData(client)
	defer clearData(client)

	trx, _ := client.BeginTransaction()

	data.Name = "Another Oska"
	trx.Table("user").ModifyOne(data)

	user1 := new(user)
	trx.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(user1)
	if user1.Name != "Another Oska" {
		t.Error("Under same transaction should receive updated data.")
	}

	trx.CommitTransaction()

	user2 := new(user)
	client.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(user2)
	if user2.Name != "Another Oska" {
		t.Error("Data should updated after transaction commit.")
	}
}

func TestUpdateRollback(t *testing.T) {
	client := getMySQLClient()
	data := setupData(client)
	defer clearData(client)

	trx, _ := client.BeginTransaction()

	data.Name = "Another Oska"
	trx.Table("user").ModifyOne(data)

	user1 := new(user)
	trx.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(user1)
	if user1.Name != "Another Oska" {
		t.Error("Under same transaction should receive updated data.")
	}

	trx.RollbackTransaction()

	user2 := new(user)
	client.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(user2)
	if user2.Name != "Oska" {
		t.Error("Data should recover to before after transaction rollback.")
	}
}

func TestRemove(t *testing.T) {
	client := getMySQLClient()
	data := setupData(client)
	defer clearData(client)

	trx, _ := client.BeginTransaction()

	trx.Table("user").DestroyOne(data)
	if err := client.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(&user{}); err != nil {
		t.Error("Data should exists with same transaction.")
	}

	trx.RollbackTransaction()

	if err := client.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(&user{}); err != nil {
		t.Error("Data should exists if transaction have been rollback.")
	}

	trx2, _ := client.BeginTransaction()
	trx2.Table("user").DestroyOne(data)
	trx2.CommitTransaction()

	if err := client.Table("user").FindOne(actions.FindOne().Where(expr.Equal("ID", data.ID))).Decode(&user{}); err == nil {
		t.Error("Data should not exists if transaction have been rollback.")
	}
}

func TestPaginateTrx(t *testing.T) {
	client := getMySQLClient()

	setupData(client)
	defer clearData(client)

	trx, _ := client.BeginTransaction()

	trx.Table("user").InsertOne(&user{ID: 2234})
	trx.Table("user").InsertOne(&user{ID: 3234})
	trx.Table("user").InsertOne(&user{ID: 4234})
	pg, _ := trx.Table("user").Paginate(actions.Paginate())

	results := make([]*user, 0)
	pg.All(&results)

	if len(results) != 4 {
		t.Error("Paginate doesn't work well in transaction mode.")
	}

	trx.RollbackTransaction()
}

func TestLockMode(t *testing.T) {
	opt := options.FindOne().SetLockMode(options.LockForRead)
	if opt.LockMode != options.LockForRead {
		t.Error("Lock mode doens't update in FindOneOptions.")
	}
}
