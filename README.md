<p align="center">
    <a href="https://github.com/si3nloong/sqlike/actions"><img src="https://github.com/si3nloong/sqlike/workflows/build/badge.svg?branch=master" alt="build status" title="build status"/></a>
    <a href="https://github.com/si3nloong/sqlike/releases"><img src="https://img.shields.io/github/v/tag/si3nloong/sqlike" alt="semver tag" title="semver tag"/></a>
    <a href="https://goreportcard.com/report/github.com/si3nloong/sqlike"><img src="https://goreportcard.com/badge/github.com/si3nloong/sqlike" alt="go report card" title="go report card"/></a>
    <a href="https://codecov.io/gh/si3nloong/sqlike"><img src="https://codecov.io/gh/si3nloong/sqlike/branch/master/graph/badge.svg" alt="coverage status" title="coverage status"/></a>
    <a href="https://github.com/si3nloong/sqlike/blob/master/LICENSE"><img src="https://img.shields.io/github/license/si3nloong/sqlike" alt="license" title="license"/></a>
</p>

# sqlike

> A golang SQL ORM which anti toxic query and focus on latest features.

## Installation

```console
go get github.com/si3nloong/sqlike
```

Fully compatible with native library `database/sql`, which mean you are allow to use `driver.Valuer` and `sql.Scanner`.

## Minimum Requirements

- **mysql 5.7** and above
- **golang 1.15** and above

## Why another ORM?

- We don't really care about _legacy support_, we want _latest feature_ that mysql and golang offer us
- We want to get rid from _toxic query_

## What do we provide apart from native package (database/sql)?

- Support `ENUM` and `SET`
- Support `UUID`
- Support `JSON`
- Support `descending index` for mysql 8.0
- Support `Spatial` with package [orb](https://github.com/paulmach/orb), such as `Point`, `LineString`
- Support `generated column` for `stored column` and `virtual column`
- Extra custom type such as `Date`, `Key`, `Boolean`
- Support `struct` on `Find`, `FindOne`, `InsertOne`, `Insert`, `ModifyOne`, `DeleteOne`, `Delete`, `DestroyOne` and `Paginate` apis
- Support `Transactions`
- Support cursor based pagination
- Support advance and complex query statement
- Support [language.Tag](https://godoc.org/golang.org/x/text/language#example-Tag--Values) and [currency.Unit](https://godoc.org/golang.org/x/text/currency#Unit)
- Support authorization plugin [Casbin](https://github.com/casbin/casbin)
- Support tracing plugin [OpenTracing](https://github.com/opentracing/opentracing-go)
- Developer friendly, (query is highly similar to native sql query)
- Support `sqldump` for backup purpose **(experiment)**

## Missing DOC?

You can refer to [examples](https://github.com/si3nloong/sqlike/tree/master/examples) folder to see what apis we offer and learn how to use those apis

## Limitation

Our main objective is anti toxic query, that why some functionality we doesn't offer out of box

- offset based pagination (but you may achieve this by using `Limit` and `Offset`)
- eager loading (we want to avoid magic function, you should handle this by your own using goroutines)
- join (eg. left join, outer join, inner join), join clause is consider as toxic query, you should alway find your record using primary key
- left wildcard search using Like is not allow (but you may use `expr.Raw` to bypass it)
- bidirectional sorting is not allow (except mysql 8.0 and above)
- currently only support `mysql` driver (postgres and sqlite yet to implement)

## General APIs

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
    ICNo      string     `sqlike:",generated_column"` // generated column generated by virtual column `Detail.ICNo`
    Name      string
    Email     string     `sqlike:",size=200,charset=latin1"` // you can set the data type length and charset with struct tag
    Address   string     `sqlike:",longtext"` // `longtext` is an alias of long text data type in mysql
    Detail    struct {
        ICNo    string `sqlike:",virtual_column=ICNo"` // virtual column
        PhoneNo string
        Age     uint
    }
    Status    UserStatus `sqlike:",enum=ACTIVE|SUSPEND"` // enum data type
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
            // `sqlike.ErrNoRows` is an alias of `sql.ErrNoRows`
            if err != sqlike.ErrNoRows {  // or you may check with sql.ErrNoRows
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
        // map into the struct of slice
        if err:= result.All(&users); err != nil {
            panic(err)
        }
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
            actions.UpdateOne().
                Where(
                    expr.Equal("ID", 100),
                ).Set(
                    expr.ColumnValue("Name", "SianLoong"),
                    expr.ColumnValue("Email", "test@gmail.com"),
                ),
            options.UpdateOne().SetDebug(true), // debug the query
        )
    }

    {
        limit := uint(10)
        pg, err := userTable.Paginate(
            ctx,
            actions.Paginate().
                OrderBy(
                    expr.Desc("CreatedAt"),
                ).
                Limit(limit + 1),
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
            length := uint(len(users))
            if length == 0 {
                break
            }
            cursor := users[length-1].ID
            if err := pg.NextCursor(ctx, cursor); err != nil {
                if err == sqlike.ErrInvalidCursor {
                    break
                }
                panic(err)
            }
            if length <= limit {
                break
            }
        }

    }
}
```

## Integrate with [OpenTracing](https://github.com/opentracing/opentracing-go)

Tracing is become very common in practice. And you may integrate your `OpenTracing` as such :

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

    itpr := opentracing.NewInterceptor(
        opentracing.WithDBInstance("sqlike"),
        opentracing.WithDBUser(cfg.User),
        opentracing.WithExec(true),
        opentracing.WithQuery(true),
    )
    client, err := sqlike.ConnectDB(
        ctx, 
        driver,
        instrumented.WrapConnector(conn, itpr),
    )
    if err != nil {
        panic(err)
    }
    defer client.Close()
}
```

Inspired by [gorm](https://github.com/jinzhu/gorm), [mongodb-go-driver](https://github.com/mongodb/mongo-go-driver) and [sqlx](https://github.com/jmoiron/sqlx).

## Special Sponsors

<p>
    <img src="https://revenuemonster.my/public/img/rm-logowhite-3x.png" alt="RevenueMonster" width="180px" style="margin:5px 10px;">
    <img src="https://asset.wetix.my/images/logo/wetix.png" alt="WeTix" width="180px" style="margin:5px 10px;">
</p>

## Big Thanks To

Thanks to these awesome companies for their support of Open Source developers ❤

[![GitHub](https://jstools.dev/img/badges/github.svg)](https://github.com/open-source)
[![NPM](https://jstools.dev/img/badges/npm.svg)](https://www.npmjs.com/)

## License

[MIT](https://github.com/si3nloong/sqlike/blob/master/LICENSE)

Copyright (c) 2019-present, SianLoong Lee
