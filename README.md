# Sequel ORM

[![Build Status](https://github.com/si3nloong/sqlike/workflows/sqlike/badge.svg?branch=master)](https://github.com/si3nloong/sqlike/actions)

```bash
go get github.com/si3nloong/sqlike
```

Fully compatible with native library `database/sql`, which mean you are allow to use `driver.Valuer` and `sql.Scanner`

### Minimum Requirements

We don't really care about _legacy support_, we want _latest feature_ that mysql and golang offer us :

- **mysql 5.7** and above
- **Golang 1.3** and above

### What we provide apart from native package (database/sql)?

- Support `JSON`
- Support `Enum`
- Support `UUID`
- Support `generated column` for `stored column` and `virtual column`
- Extra custom type such as `Date`, `Key`, `Boolean`
- Support `struct` on `Find`, `FindOne`, `InsertOne`, `Insert`, `ModifyOne`, `DeleteOne`, `Delete`, `DestroyOne` and `Paginate` apis
- Support [language.Tag](https://godoc.org/golang.org/x/text/language#example-Tag--Values) and [currency.Unit](https://godoc.org/golang.org/x/text/currency#Unit)
- Support third-party plugin [Casbin](https://github.com/casbin/casbin)
- Prevent toxic query with `Strict Mode` **(in progress)**
- Support query filtering **(in progress)**

### Missing DOC?

You can refer to [examples](https://github.com/si3nloong/sqlike/tree/master/examples) folder to see what apis we offer and learn how to use those apis

```go
import (
    "time"
    "github.com/si3nloong/sqlike/sqlike/actions"
    "github.com/si3nloong/sqlike/sqlike"
    "github.com/si3nloong/sqlike/sqlike/options"
    "github.com/si3nloong/sqlike/sql/expr"
    "github.com/google/uuid"
)

// UserStatus :
type UserStatus string

const (
    UserStatusActive  UserStatus = "ACTIVE"
    UserStatusSuspend UserStatus = "SUSPEND"
)

type User struct {
    ID        uuid.UUID
    ICNo      string     `sqlike:",generated_column"`
    Name      string
    Email     string     `sqlike:",size=200"`
    Address   string     `sqlike:",longtext"`
    Detail    struct {
        ICNo    string `sqlike:",virtual_column=ICNo"`
        PhoneNo string
        Age     uint
    }
    Status    UserStatus `sqlike:",charset=latin1,enum=ACTIVE|SUSPEND"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func newUser() (user User) {
    now := time.Now()
    user.ID = uuid.New()
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
        SetPort("3306"),
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
        if _, err := userTable.Insert(&users); err != nil {
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
        result, err := userTable.Find(
            actions.Find().Where(
                expr.Equal("ID", result.ID),
            ),
        )
        if err != nil {
            panic(err)
        }
        result.All(&users) // map into the struct of slice
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
                expr.ColumnValue("Name", "SianLoong"),
                expr.ColumnValue("Email", "test@gmail.com"),
            ),
            options.UpdateOne().SetDebug(true), // debug the query
        )
    }

    {
        pg, err := userTable.Paginate(
            actions.Paginate().
                OrderBy(
                    expr.Desc("CreatedAt"),
                ),
             options.Paginate().SetDebug(true),
        )
        if err != nil {
            panic(err)
        }

        for {
            var users []User
            if err := pg.All(&users); err != nil {
                panic(err)
            }
            if len(users) == 0 {
                break
            }
            cursor := users[len(users)-1].ID
            if err := pg.NextPage(cursor); err != nil {
                if err == sqlike.ErrInvalidCursor {
                    break
                }
                panic(err)
            }
        }

    }
}
```
