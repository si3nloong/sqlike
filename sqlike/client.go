package sqlike

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/si3nloong/sqlike/sql/charset"
	"github.com/si3nloong/sqlike/sql/codec"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/logs"
)

// Client :
type Client struct {
	driverName string
	version    *semver.Version
	*sql.DB
	pk      string
	logger  logs.Logger
	charSet charset.Code
	collate string
	dialect sqldialect.Dialect
}

func newClient(driver string, db *sql.DB, dialect sqldialect.Dialect, code charset.Code, collate string) (*Client, error) {
	driver = strings.TrimSpace(strings.ToLower(driver))
	client := &Client{
		driverName: driver,
		DB:         db,
		dialect:    dialect,
		charSet:    code,
		collate:    collate,
	}
	client.pk = "$Key"
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

// SetPrimaryKey :
func (c *Client) SetPrimaryKey(pk string) *Client {
	c.pk = pk
	return c
}

// Version :
func (c *Client) Version() (version *semver.Version) {
	if c.version == nil {
		c.version = c.getVersion()
	}
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
		driverName: c.driverName,
		name:       name,
		pk:         c.pk,
		client:     c,
		dialect:    c.dialect,
		driver:     c.DB,
		logger:     c.logger,
		registry:   codec.DefaultRegistry,
	}
}

func (c *Client) getVersion() (version *semver.Version) {
	var (
		ver string
		err error
	)
	stmt := c.dialect.GetVersion()
	err = sqldriver.QueryRowContext(
		context.Background(),
		c.DB,
		stmt,
		c.logger,
	).Scan(&ver)
	if err != nil {
		panic(err)
	}
	version, err = semver.NewVersion(ver)
	if err != nil {
		panic(err)
	}
	return
}
