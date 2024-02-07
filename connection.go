package sqlike

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"

	"github.com/si3nloong/sqlike/v2/options"
	"github.com/si3nloong/sqlike/v2/sql/charset"
	"github.com/si3nloong/sqlike/v2/sql/dialect"
	"github.com/si3nloong/sqlike/v2/sql/dialect/mysql"
)

func init() {
	dialect.RegisterDialect("mysql", mysql.New())
}

// Open : open connection to sql server with connection string
func Open(ctx context.Context, driver string, opt *options.ConnectOptions) (client *Client, err error) {
	if opt == nil {
		return nil, errors.New("sqlike: invalid connection options <nil>")
	}
	var database *sql.DB
	dlct := dialect.GetDialectByDriver(driver)
	connStr := dlct.Connect(opt)
	log.Println("Connecting to :", connStr)
	database, err = sql.Open(driver, connStr)
	if err != nil {
		return
	}
	client, err = newClient(ctx, driver, database, dlct, opt.Charset, opt.Collate)
	return
}

// MustOpen : must open will panic if it cannot establish a connection to sql server
func MustOpen(ctx context.Context, driver string, opt *options.ConnectOptions) *Client {
	client, err := Open(ctx, driver, opt)
	if err != nil {
		panic(err)
	}
	return client
}

// Connect : connect and ping the sql server, throw error when unable to ping
func Connect(ctx context.Context, driver string, opt *options.ConnectOptions) (client *Client, err error) {
	client, err = Open(ctx, driver, opt)
	if err != nil {
		return
	}
	err = client.PingContext(ctx)
	if err != nil {
		client.Close()
		return
	}
	return
}

// MustConnect will panic if cannot connect to sql server
func MustConnect(ctx context.Context, driver string, opt *options.ConnectOptions) *Client {
	conn, err := Connect(ctx, driver, opt)
	if err != nil {
		panic(err)
	}
	return conn
}

// ConnectDB :
func ConnectDB(ctx context.Context, driver string, conn driver.Connector) (*Client, error) {
	database := sql.OpenDB(conn)
	dlct := dialect.GetDialectByDriver(driver)
	client, err := newClient(
		ctx,
		driver,
		database,
		dlct,
		charset.DefaultCharset,
		charset.DefaultCollation,
	)
	if err != nil {
		return nil, err
	}
	if err := client.PingContext(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

// MustConnectDB :
func MustConnectDB(ctx context.Context, driver string, conn driver.Connector) *Client {
	client, err := ConnectDB(ctx, driver, conn)
	if err != nil {
		panic(err)
	}
	return client
}
