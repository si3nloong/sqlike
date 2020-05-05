# Sequel ORM

[![Build Status](https://github.com/si3nloong/sqlike/workflows/build/badge.svg?branch=master)](https://github.com/si3nloong/sqlike/actions)

```bash
go get github.com/si3nloong/sqlike
```

Fully compatible with native library `database/sql`, which mean you are allow to use `driver.Valuer` and `sql.Scanner`.

### Minimum Requirements

- **mysql 5.7** and above
- **Golang 1.13** and above

### Why another ORM?

- We don't really care about _legacy support_, we want _latest feature_ that mysql and golang offer us
- We want to get rid from _toxic query_

### What we provide apart from native package (database/sql)?

- Support `ENUM` and `SET`
- Support `UUID`
- Support `JSON`
- Support `descending index` for mysql 8.0
- Support `Spatial` with package [orb](https://github.com/paulmach/orb), such as `Point`, `LineString`
- Support `generated column` for `stored column` and `virtual column`
- Extra custom type such as `Date`, `Key`, `Boolean`
- Support `struct` on `Find`, `FindOne`, `InsertOne`, `Insert`, `ModifyOne`, `DeleteOne`, `Delete`, `DestroyOne` and `Paginate` apis
- Support `Transactions`
- Support cursor based pagination.
- Support advance and complex query statement
- Support [language.Tag](https://godoc.org/golang.org/x/text/language#example-Tag--Values) and [currency.Unit](https://godoc.org/golang.org/x/text/currency#Unit)
- Support authorization plugin [Casbin](https://github.com/casbin/casbin)
- Support tracing plugin [OpenTracing](https://github.com/opentracing/opentracing-go)
- Support `sqldump` for backup purpose **(in progress)**
- Developer friendly, (query is highly similar to native sql query)
- Prevent toxic query with `Strict Mode` **(upcoming)**

### Missing DOC?

You can refer to [examples](https://github.com/si3nloong/sqlike/tree/master/examples) folder to see what apis we offer and learn how to use those apis

### Limitation

Our main objective is anti toxic query, that why some functionality we doesn't offer

- eager loading (we want to avoid magic function, you should handle this by your own using goroutines)
- join (eg. left join, outer join, inner join)
- left wildcard search using Like is not allow (you may use Raw to bypass)
- bidirectional sorting is not allow (except mysql 8.0)
- currently only support `mysql` driver

### General APIs

```go
import (
    "time"
    "github.com/si3nloong/sqlike/sqlike/actions"
    "github.com/si3nloong/sqlike/sqlike"
    "github.com/si3nloong/sqlike/sqlike/options"
    "github.com/si3nloong/sqlike/sql/expr"
    "github.com/google/uuid"
    "context"
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
    ctx := context.Background()
    client := sqlike.MustConnect(
        ctx,
        "mysql",
        options.Connect().
        SetUsername("root").
        SetPassword("").
        SetHost("localhost").
        SetPort("3306"),
    )

    client.SetPrimaryKey("ID") // Change default primary key name
    version := client.Version() // Display driver version
    dbs, _ := client.ListDatabases(ctx) // List databases

    userTable := client.Database("sqlike").Table("User")

    // Drop Table
    userTable.Drop(ctx)

    // Migrate Table
    userTable.Migrate(ctx, User{})

    // Truncate Table
    userTable.Truncate(ctx)

    // Insert one record
    {
        user := newUser()
        if _, err := userTable.InsertOne(ctx, &user); err != nil {
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
        if _, err := userTable.Insert(ctx, &users); err != nil {
            panic(err)
        }
    }

    // Find one record
    {
        user := User{}
        err := userTable.FindOne(ctx, nil).Decode(&user)
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
            ctx,
            actions.Find().
                Where(
                    expr.Equal("ID", result.ID),
                ).
                OrderBy(
                    expr.Desc("UpdatedAt"),
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
        if err := userTable.ModifyOne(ctx, &user); err != nil {
            panic(err)
        }
    }

    // Update one record with selected fields
    {
        userTable.UpdateOne(
            ctx,
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
            ctx,
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
            if err := pg.NextCursor(ctx, cursor); err != nil {
                if err == sqlike.ErrInvalidCursor {
                    break
                }
                panic(err)
            }
        }

    }
}
```

### Integrate with [OpenTracing](https://github.com/opentracing/opentracing-go)

Tracing is become more and more common in practice. And you may integrate your `OpenTracing` as such :

```go
import (
    "context"
    "github.com/go-sql-driver/mysql"
    "github.com/si3nloong/sqlike/plugin/opentracing"
    "github.com/si3nloong/sqlike/sql/instrumented"

    "github.com/si3nloong/sqlike/sqlike"
)

func main() {
    ctx := context.Background()
    driver := "mysql"
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
    client, err := sqlike.ConnectDB(
        ctx, driver,
        instrumented.WrapConnector(conn, itpr),
    )
    if err != nil {
        panic(err)
    }
    defer client.Close()
}
```

Inspired by [gorm](https://github.com/jinzhu/gorm), [mongodb-go-driver](https://github.com/mongodb/mongo-go-driver) and [sqlx](https://github.com/jmoiron/sqlx).
