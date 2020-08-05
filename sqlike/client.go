package sqlike

import (
	"context"
	"database/sql"
	"strings"

	semver "github.com/Masterminds/semver/v3"
	"github.com/si3nloong/sqlike/reflext"
	"github.com/si3nloong/sqlike/sql/charset"
	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
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

// Client : sqlike client is a client embedded with *sql.DB, so you may use any apis of *sql.DB
type Client struct {
	*DriverInfo
	*sql.DB
	pk      string
	logger  logs.Logger
	cache   reflext.StructMapper
	codec   codec.Codecer
	dialect dialect.Dialect
}

// newClient : create a new client struct by providing driver, *sql.DB, dialect etc
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
	client.cache = reflext.DefaultMapper
	client.codec = codec.DefaultRegistry
	client.version = client.getVersion(ctx)
	return client, nil
}

// SetLogger : this is to set the logger for debugging, it will panic if the logger input is nil
func (c *Client) SetLogger(logger logs.Logger) *Client {
	if logger == nil {
		panic("logger cannot be nil")
	}
	c.logger = logger
	return c
}

// SetPrimaryKey : this will set a default primary key for subsequent operation such as Insert, InsertOne, ModifyOne
func (c *Client) SetPrimaryKey(pk string) *Client {
	c.pk = pk
	return c
}

// SetCodec : Codec is a component which handling the :
// 1. encoding between input data and driver.Valuer
// 2. decoding between output data and sql.Scanner
func (c *Client) SetCodec(cdc codec.Codecer) *Client {
	c.codec = cdc
	return c
}

// SetStructMapper : StructMapper is a mapper to reflect a struct on runtime and provide struct info
func (c *Client) SetStructMapper(mapper reflext.StructMapper) *Client {
	c.cache = mapper
	return c
}

// CreateDatabase : create database with name
func (c *Client) CreateDatabase(ctx context.Context, name string) error {
	return c.createDB(ctx, name, true)
}

// DropDatabase : drop the selected database
func (c *Client) DropDatabase(ctx context.Context, name string) error {
	return c.dropDB(ctx, name, true)
}

// ListDatabases : list all the database on current connection
func (c *Client) ListDatabases(ctx context.Context) ([]string, error) {
	stmt := sqlstmt.AcquireStmt(c.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	c.dialect.GetDatabases(stmt)
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

// Database : this api will execute `USE database`, which will point your current connection to selected database
func (c *Client) Database(name string) *Database {
	stmt := sqlstmt.AcquireStmt(c.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	c.dialect.UseDatabase(stmt, name)
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

// getVersion is a internal function to get sql driver's version
func (c *Client) getVersion(ctx context.Context) (version *semver.Version) {
	var (
		ver string
		err error
	)
	stmt := sqlstmt.AcquireStmt(c.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	c.dialect.GetVersion(stmt)
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

// createDB is a internal function for create a database
func (c *Client) createDB(ctx context.Context, name string, checkExists bool) error {
	stmt := sqlstmt.AcquireStmt(c.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	c.dialect.CreateDatabase(stmt, name, checkExists)
	_, err := driver.Execute(
		ctx,
		c.DB,
		stmt,
		c.logger,
	)
	return err
}

// dropDB is a internal function for drop a database
func (c *Client) dropDB(ctx context.Context, name string, checkExists bool) error {
	stmt := sqlstmt.AcquireStmt(c.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	c.dialect.DropDatabase(stmt, name, checkExists)
	_, err := driver.Execute(
		ctx,
		c.DB,
		stmt,
		c.logger,
	)
	return err
}
