package sqlike

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/si3nloong/sqlike/sql/charset"
	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/logs"
)

// DriverInfo :
type DriverInfo struct {
	driverName string
	version    *semver.Version
	charSet    charset.Code
	collate    string
}

// DriverName :
func (d *DriverInfo) DriverName() string {
	return d.driverName
}

// Version :
func (d *DriverInfo) Version() *semver.Version {
	return d.version
}

// Charset :
func (d *DriverInfo) Charset() charset.Code {
	return d.charSet
}

// Collate :
func (d *DriverInfo) Collate() string {
	return d.collate
}

// Client :
type Client struct {
	*DriverInfo
	*sql.DB
	pk      string
	logger  logs.Logger
	codec   codec.Codecer
	dialect dialect.Dialect
}

func newClient(ctx context.Context, driver string, db *sql.DB, dialect dialect.Dialect, code charset.Code, collate string) (*Client, error) {
	driver = strings.TrimSpace(strings.ToLower(driver))
	client := &Client{
		DB:      db,
		dialect: dialect,
	}
	client.pk = "$Key"
	client.DriverInfo = new(DriverInfo)
	client.driverName = driver
	client.charSet = code
	client.collate = collate
	client.codec = codec.DefaultRegistry
	client.version = client.getVersion(ctx)
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

// SetCodec :
func (c *Client) SetCodec(cdc codec.Codecer) *Client {
	c.codec = cdc
	return c
}

// CreateDatabase :
func (c *Client) CreateDatabase(ctx context.Context, name string) error {
	return c.createDB(ctx, name, true)
}

// DropDatabase :
func (c *Client) DropDatabase(ctx context.Context, name string) error {
	return c.dropDB(ctx, name, true)
}

// ListDatabases :
func (c *Client) ListDatabases(ctx context.Context) ([]string, error) {
	stmt := c.dialect.GetDatabases()
	rows, err := driver.Query(
		ctx,
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
	stmt := c.dialect.UseDatabase(name)
	if _, err := driver.Execute(context.Background(), c.DB, stmt, c.logger); err != nil {
		panic(err)
	}
	return &Database{
		driverName: c.driverName,
		name:       name,
		pk:         c.pk,
		client:     c,
		dialect:    c.dialect,
		driver:     c.DB,
		logger:     c.logger,
		codec:      c.codec,
	}
}

func (c *Client) getVersion(ctx context.Context) (version *semver.Version) {
	var (
		ver string
		err error
	)
	stmt := c.dialect.GetVersion()
	err = driver.QueryRowContext(
		ctx,
		c.DB,
		stmt,
		c.logger,
	).Scan(&ver)
	if err != nil {
		panic(err)
	}
	paths := strings.Split(ver, "-")
	version, err = semver.NewVersion(paths[0])
	if err != nil {
		panic(err)
	}
	return
}

func (c *Client) createDB(ctx context.Context, name string, checkExists bool) error {
	stmt := c.dialect.CreateDatabase(name, checkExists)
	_, err := driver.Execute(
		ctx,
		c,
		stmt,
		c.logger,
	)
	return err
}

func (c *Client) dropDB(ctx context.Context, name string, checkExists bool) error {
	stmt := c.dialect.DropDatabase(name, checkExists)
	_, err := driver.Execute(
		ctx,
		c,
		stmt,
		c.logger,
	)
	return err
}
