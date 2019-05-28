package sqlike

import "database/sql"

// Open : connect to sql server with connection string
func Open(driver, connStr string) (client *Client, err error) {
	var conn *sql.DB
	conn, err = sql.Open(driver, connStr)
	if err != nil {
		return
	}
	client, err = newClient(driver, conn)
	return
}

// Connect :
func Connect(driver, connStr string) (client *Client, err error) {
	client, err = Open(driver, connStr)
	err = client.Ping()
	if err != nil {
		client.Close()
		return
	}
	return
}

// MustConnect will panic if cannot connect to sql server
func MustConnect(driver, connStr string) *Client {
	conn, err := Connect(driver, connStr)
	if err != nil {
		panic(err)
	}
	return conn
}
