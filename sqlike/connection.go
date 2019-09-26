package sqlike

import (
	"database/sql"
	"log"

	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sqlike/options"
)

// Open : connect to sql server with connection string
func Open(driver string, opt *options.ConnectOptions) (client *Client, err error) {
	if opt == nil {
		panic("sqlike: invalid connection options <nil>")
	}
	var conn *sql.DB
	dialect := sqldialect.GetDialectByDriver(driver)
	connStr := dialect.Connect(opt)
	log.Println("Connect to :", connStr)
	conn, err = sql.Open(driver, connStr)
	if err != nil {
		return
	}
	client, err = newClient(driver, conn, dialect, opt.Charset, opt.Collate)
	return
}

// Connect :
func Connect(driver string, opt *options.ConnectOptions) (client *Client, err error) {
	client, err = Open(driver, opt)
	err = client.Ping()
	if err != nil {
		client.Close()
		return
	}
	return
}

// MustConnect will panic if cannot connect to sql server
func MustConnect(driver string, opt *options.ConnectOptions) *Client {
	conn, err := Connect(driver, opt)
	if err != nil {
		panic(err)
	}
	return conn
}
