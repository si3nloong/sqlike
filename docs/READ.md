# Read

### Retrieving single record

```go
import (
    "github.com/si3nloong/sqlike/sql/expr"
    "github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
)

var user User{
    Name string
    Age  int
}

db.Table("users").
    FindOne(
        context.Background(),
        actions.FindOne().
            Where(
                expr.Equal("ID", "123"),
            ),
    ).
    Decode(&user)
```

### Retrieving multiple records

```go
import (
    "github.com/si3nloong/sqlike/sql/expr"
    "github.com/si3nloong/sqlike/sqlike/actions"
	"github.com/si3nloong/sqlike/sqlike/options"
)

result, err := db.Table("users").
    Find(
        context.Background(),
        actions.Find(),
    )

users := []User{}
if err := result.All(&users); err != nil {
    panic(err)
}
```
