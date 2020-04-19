package sqlike

import (
	"context"
	"database/sql"
	"errors"
	"log"

	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// Open : open connection to sql server with connection string
func Open(ctx context.Context, driver string, opt *options.ConnectOptions) (client *Client, err error) {
	if opt == nil {
		return nil, errors.New("sqlike: invalid connection options <nil>")
	}
	var conn *sql.DB
	dialect := sqldialect.GetDialectByDriver(driver)
	connStr := dialect.Connect(opt)
	log.Println("Connect to :", connStr)
	conn, err = sql.Open(driver, connStr)
	if err != nil {
		return
	}
	client, err = newClient(ctx, driver, conn, dialect, opt.Charset, opt.Collate)
	return
}

// MustOpen :
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
