package sqlike

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/si3nloong/sqlike/sql/codec"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/si3nloong/sqlike/sqlike/logs"
	"github.com/si3nloong/sqlike/sqlike/options"
	"gopkg.in/yaml.v3"
)

// Database :
type Database struct {
	driverName string
	name       string
	pk         string
	client     *Client
	driver     sqldriver.Driver
	dialect    sqldialect.Dialect
	registry   *codec.Registry
	logger     logs.Logger
}

// Create :
func (db *Database) Create() error {
	return db.createDB(true)
}

// CreateIfNotExists :
func (db *Database) CreateIfNotExists() error {
	return db.createDB(false)
}

// Drop :
func (db *Database) Drop() error {
	return db.dropDB(false)
}

// DropIfExists :
func (db *Database) DropIfExists() error {
	return db.dropDB(true)
}

// Table :
func (db *Database) Table(name string) *Table {
	return &Table{
		dbName:   db.name,
		name:     name,
		pk:       db.pk,
		client:   db.client,
		driver:   db.driver,
		dialect:  db.dialect,
		registry: db.registry,
		logger:   db.logger,
	}
}

// BeginTransaction :
func (db *Database) BeginTransaction() (*Transaction, error) {
	return db.beginTrans(context.Background(), nil)
}

func (db *Database) beginTrans(ctx context.Context, opt *sql.TxOptions) (*Transaction, error) {
	tx, err := db.client.BeginTx(ctx, opt)
	if err != nil {
		return nil, err
	}
	return &Transaction{
		dbName:   db.name,
		pk:       db.pk,
		client:   db.client,
		context:  ctx,
		driver:   tx,
		dialect:  db.dialect,
		logger:   db.logger,
		registry: db.registry,
	}, nil
}

type txCallback func(ctx SessionContext) error

// RunInTransaction :
func (db *Database) RunInTransaction(cb txCallback, opts ...*options.TransactionOptions) error {
	opt := new(options.TransactionOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	duration := 60 * time.Second
	if opt.Duration.Seconds() > 0 {
		duration = opt.Duration
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	tx, err := db.beginTrans(ctx, &sql.TxOptions{
		Isolation: opt.IsolationLevel,
		ReadOnly:  opt.ReadOnly,
	})
	if err != nil {
		return err
	}
	defer tx.RollbackTransaction()
	if err := cb(tx); err != nil {
		return err
	}
	return tx.CommitTransaction()
}

type indexDefinition struct {
	Indexes []struct {
		Table   string `yaml:"table"`
		Name    string `yaml:"name"`
		Type    string `yaml:"type"`
		Columns []struct {
			Name      string `yaml:"name"`
			Direction string `yaml:"direction"`
		} `yaml:"columns"`
	} `yaml:"indexes"`
}

// BuildIndexes :
func (db *Database) BuildIndexes(filepath ...string) error {
	var id indexDefinition
	pwd, _ := os.Getwd()
	file := pwd + "/index.yaml"
	if len(filepath) > 0 {
		file = filepath[0]
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, &id); err != nil {
		return err
	}

	for _, idx := range id.Indexes {
		columns := make([]indexes.Column, 0, len(idx.Columns))
		for _, col := range idx.Columns {
			dir := indexes.Ascending
			col.Direction = strings.TrimSpace(strings.ToLower(col.Direction))
			if col.Direction == "desc" || col.Direction == "descending" {
				dir = indexes.Descending
			}
			columns = append(columns, indexes.Column{
				Name:      col.Name,
				Direction: dir,
			})
		}

		index := indexes.Index{
			Name:    idx.Name,
			Type:    parseIndexType(idx.Type),
			Columns: columns,
		}

		if exists, err := isIndexExists(
			db.name,
			idx.Table,
			index.GetName(),
			db.driver,
			db.dialect,
			db.logger,
		); err != nil {
			return err
		} else if exists {
			continue
		}

		iv := db.Table(idx.Table).Indexes()
		if err := iv.CreateOne(index); err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) createDB(exists bool) error {
	stmt := db.dialect.CreateDatabase(db.name, exists)
	_, err := sqldriver.Execute(
		context.Background(),
		db.driver,
		stmt,
		db.logger,
	)
	return err
}

func (db *Database) dropDB(exists bool) error {
	stmt := db.dialect.DropDatabase(db.name, exists)
	_, err := sqldriver.Execute(
		context.Background(),
		db.driver,
		stmt,
		db.logger,
	)
	return err
}

func parseIndexType(name string) (idxType indexes.Type) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		idxType = indexes.BTree
		return
	}

	switch name {
	case "spatial":
		idxType = indexes.Spatial
	case "unique":
		idxType = indexes.Unique
	case "btree":
		idxType = indexes.BTree
	case "fulltext":
		idxType = indexes.FullText
	case "primary":
		idxType = indexes.Primary
	default:
		panic(fmt.Errorf("invalid index type %q", name))
	}
	return
}
