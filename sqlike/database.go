package sqlike

import (
	"bitbucket.org/SianLoong/sqlike/core/codec"
	sqlcore "bitbucket.org/SianLoong/sqlike/sqlike/sql/core"
	sqldriver "bitbucket.org/SianLoong/sqlike/sqlike/sql/driver"
)

// Database :
type Database struct {
	name     string
	client   *Client
	driver   sqldriver.Driver
	dialect  sqlcore.Dialect
	registry *codec.Registry
	logger   Logger
}

// Table :
func (db *Database) Table(name string) *Table {
	return &Table{
		dbName:   db.name,
		name:     name,
		client:   db.client,
		driver:   db.driver,
		dialect:  db.dialect,
		registry: db.registry,
		logger:   db.logger,
	}
}
