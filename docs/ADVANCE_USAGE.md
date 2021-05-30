# Advance Usage

### Integrate with [OpenTracing](https://github.com/opentracing/opentracing-go)

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
