package sqlike

import (
	"context"
	"database/sql"

	"github.com/blang/semver"
	"github.com/si3nloong/sqlike/sqlike/logs"
	"github.com/si3nloong/sqlike/sqlike/sql/codec"
	sqlcore "github.com/si3nloong/sqlike/sqlike/sql/core"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
)

// Client :
type Client struct {
	driverName string
	version    semver.Version
	*sql.DB
	logger  logs.Logger
	dialect sqlcore.Dialect
}

func newClient(driver string, db *sql.DB, dialect sqlcore.Dialect) (*Client, error) {
	client := &Client{
		driverName: driver,
		DB:         db,
		dialect:    dialect,
	}
	client.version = client.getVersion()
	return client, nil
}

// SetLogger : this is to set the logger for debugging, it will panic
// if the logger input is nil
func (c *Client) SetLogger(logger logs.Logger) *Client {
	if logger == nil {
		panic("logger cannot be nil")
	}
	c.logger = logger
	return c
}

// Version :
func (c *Client) Version() (version semver.Version) {
	version = c.version
	return
}

// ListDatabases :
func (c Client) ListDatabases() ([]string, error) {
	stmt := c.dialect.GetDatabases()
	rows, err := sqldriver.Query(
		context.Background(),
		c.DB,
		stmt,
		c.logger,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dbs := make([]string, 0)
	for i := 0; rows.Next(); i++ {
		dbs = append(dbs, "")
		if err := rows.Scan(&dbs[i]); err != nil {
			return nil, err
		}
	}
	return dbs, nil
}

// Database :
func (c *Client) Database(name string) *Database {
	return &Database{
		name:     name,
		client:   c,
		dialect:  c.dialect,
		driver:   c.DB,
		logger:   c.logger,
		registry: codec.DefaultRegistry,
	}
}

func (c *Client) getVersion() (version semver.Version) {
	stmt := c.dialect.GetVersion()
	var ver string
	sqldriver.QueryRowContext(
		context.Background(),
		c.DB,
		stmt,
		c.logger,
	).Scan(&ver)
	version, _ = semver.Parse(ver)
	return
}
