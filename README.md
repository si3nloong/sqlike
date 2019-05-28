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
2. Support `stored column` and `virtual column`
3. Support `struct` on `Find`, `FindOne`, `InsertOne`, `InsertMany` and `ModifyOne` apis

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
    ID        uuid.UUID  `sqlike:"$Key"`
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
        SetUsername(`root`).
        SetPassword(``).
        SetHost(`localhost`).
        SetPort(`3306`).
        SetDatabase(`sqlike`),
    )

    version := client.Version() // Display driver version
    dbs, _ := client.ListDatabases() // List databases

    userTable := client.Database("sqlike").Table("User")
    userTable.Drop() // Drop Table
    userTable.Migrate(new(User)) // Migration

    user := newUser()
    if _, err := userTable.InsertOne(&user); err != nil {
        panic(err)
    }

    users := [...]User{
        newUser(),
        newUser(),
        newUser(),
    }
    if _, err := userTable.InsertMany(&users); err != nil {
        panic(err)
    }

    user.Name = `🤖 Hello World!`
    if err := userTable.ModifyOne(&user); err != nil {
        panic(err)
    }

    result := User{}
    if err := userTable.FindOne(nil).Decode(&result); err != nil {
        panic(err)
    }


    cursor, err := userTable.Find(
        actions.Find().
        Where(expr.Equal("$Key", result.ID)),
    )
}
```
