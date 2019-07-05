# Sequel ORM

```bash
go get github.com/si3nloong/sqlike
```

Fully compatible with native library `database/sql`, which mean you are allow to use `driver.Valuer` and `sql.Scanner`

### Minimum Requirements

We don't really care about legacy support, we want latest feature that mysql and golang offer us :

1. mysql 5.7 and above
2. Golang 1.10 and above

### What we provide apart from native package (database/sql)?

1. Support `JSON`
2. Support `Enum`
3. Support `UUID`
4. Support `stored column` and `virtual column`
5. Extra type such as `Date`, `Timestamp`, `Key` and `GeoPoint`
6. Support `struct` on `Find`, `FindOne`, `InsertOne`, `InsertMany` and `ModifyOne` apis

```go
import (
    "time"
    "github.com/si3nloong/sqlike/sqlike/actions"
    "github.com/si3nloong/sqlike/sqlike"
    "github.com/si3nloong/sqlike/sqlike/options"
    "github.com/si3nloong/sqlike/sqlike/sql/expr"
    uuid "github.com/satori/go.uuid"
)

// UserStatus :
type UserStatus string

const (
    UserStatusActive  UserStatus = "ACTIVE"
    UserStatusSuspend UserStatus = "SUSPEND"
)

type User struct {
    ID        uuid.UUID 
    Name      string
    Email     string     `sqlike:",size:200"`
    Address   string     `sqlike:",longtext"`
    Status    UserStatus `sqlike:",enum:ACTIVE|SUSPEND"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func newUser() (user User) {
    now := time.Now()
    user.ID = uuid.NewV1()
    user.CreatedAt = now
    user.UpdatedAt = now
    return
}

func main() {
    client := sqlike.MustConnect("mysql",
        options.Connect().
        SetUsername("root").
        SetPassword("").
        SetHost("localhost").
        SetPort("3306").
        SetDatabase("sqlike"),
    )

    client.SetPrimaryKey("ID") // Change default primary key name
    version := client.Version() // Display driver version
    dbs, _ := client.ListDatabases() // List databases

    userTable := client.Database("sqlike").Table("User")

    // Drop Table
    userTable.Drop() 

    // Migrate Table
    userTable.Migrate(User{})

    // Truncate Table
    userTable.Truncate()

    // Insert one record
    {
        user := newUser()
        if _, err := userTable.InsertOne(&user); err != nil {
            panic(err)
        }
    }

    // Insert multiple record
    {
        users := [...]User{
            newUser(),
            newUser(),
            newUser(),
        }
        if _, err := userTable.InsertMany(&users); err != nil {
            panic(err)
        }
    }

    // Find one record
    {
        user := User{}
        err := userTable.FindOne(nil).Decode(&user)
        if err != nil {
            if err != sqlike.ErrNoRows {
                panic(err)
            }
            // record not exist
        }
    }

    // Find multiple records
    {
        users := make([]User, 0)
        cursor, err := userTable.Find(
            actions.Find().Where(
                expr.Equal("ID", result.ID),
            ),
        )
        if err != nil {
            panic(err)
        }
        cursor.All(&users) // map into the struct of slice
    }

    // Update one record with all fields of struct
    {
        user.Name = `🤖 Hello World!`
        if err := userTable.ModifyOne(&user); err != nil {
            panic(err)
        }
    }

    // Update one record with selected fields
    {
        userTable.UpdateOne(
            actions.UpdateOne().Where(
                expr.Equal("ID", 100),
            ).Set(
                expr.Field("Name", "SianLoong"),
                expr.Field("Email", "test@gmail.com"),
            ),
            options.UpdateOne().SetDebug(true), // debug the query
        )
    }
}
```
