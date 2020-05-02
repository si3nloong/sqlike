package sqlike

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/si3nloong/sqlike/sql/codec"
	"github.com/si3nloong/sqlike/sql/dialect"
	"github.com/si3nloong/sqlike/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/si3nloong/sqlike/sqlike/logs"
	"github.com/si3nloong/sqlike/sqlike/options"
	"gopkg.in/yaml.v3"
)

type txCallback func(ctx SessionContext) error

// Database :
type Database struct {
	driverName string
	name       string
	pk         string
	client     *Client
	driver     driver.Driver
	dialect    dialect.Dialect
	codec      codec.Codecer
	logger     logs.Logger
}

// Table :
func (db *Database) Table(name string) *Table {
	return &Table{
		dbName:  db.name,
		name:    name,
		pk:      db.pk,
		client:  db.client,
		driver:  db.driver,
		dialect: db.dialect,
		codec:   db.codec,
		logger:  db.logger,
	}
}

func (db *Database) QueryStmt(ctx context.Context, query interface{}) (*Result, error) {
	if query == nil {
		return nil, errors.New("empty query statement")
	}
	stmt := sqlstmt.AcquireStmt(db.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	if err := db.dialect.SelectStmt(stmt, query); err != nil {
		return nil, err
	}
	rows, err := driver.Query(
		ctx,
		db.driver,
		stmt,
		getLogger(db.logger, true),
	)
	if err != nil {
		return nil, err
	}
	rslt := new(Result)
	rslt.codec = db.codec
	rslt.rows = rows
	rslt.columnTypes, rslt.err = rows.ColumnTypes()
	if rslt.err != nil {
		defer rslt.rows.Close()
	}
	for _, col := range rslt.columnTypes {
		rslt.columns = append(rslt.columns, col.Name())
	}
	return rslt, rslt.err
}

// BeginTransaction :
func (db *Database) BeginTransaction(ctx context.Context, opts ...*sql.TxOptions) (*Transaction, error) {
	opt := &sql.TxOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	return db.beginTrans(ctx, opt)
}

func (db *Database) beginTrans(ctx context.Context, opt *sql.TxOptions) (*Transaction, error) {
	tx, err := db.client.BeginTx(ctx, opt)
	if err != nil {
		return nil, err
	}
	return &Transaction{
		Context: ctx,
		dbName:  db.name,
		pk:      db.pk,
		client:  db.client,
		driver:  tx,
		dialect: db.dialect,
		logger:  db.logger,
		codec:   db.codec,
	}, nil
}

// RunInTransaction :
func (db *Database) RunInTransaction(ctx context.Context, cb txCallback, opts ...*options.TransactionOptions) error {
	opt := new(options.TransactionOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	duration := 60 * time.Second
	if opt.Duration.Seconds() > 0 {
		duration = opt.Duration
	}
	c, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	tx, err := db.beginTrans(c, &sql.TxOptions{
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
func (db *Database) BuildIndexes(ctx context.Context, paths ...string) error {
	var (
		path string
		err  error
		fi   os.FileInfo
	)
	if len(paths) > 0 {
		path = paths[0]
		fi, err = os.Stat(path)
		if err != nil {
			return err
		}
	} else {
		pwd, _ := os.Getwd()
		files := []string{pwd + "/index.yml", pwd + "/index.yaml"}
		for _, f := range files {
			fi, err = os.Stat(f)
			if !os.IsNotExist(err) {
				path = f
				break
			}
		}
		if err != nil {
			return err
		}
	}

	switch v := fi.Mode(); {
	case v.IsDir():
		if err := filepath.Walk(path, func(fp string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			ext := filepath.Ext(info.Name())
			// only interested on yaml and yml files
			if ext != ".yaml" && ext != ".yml" {
				return nil
			}

			return db.readAndBuildIndexes(ctx, fp)
		}); err != nil {
			return err
		}

	case v.IsRegular():
		if err := db.readAndBuildIndexes(ctx, path); err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) readAndBuildIndexes(ctx context.Context, path string) error {
	var id indexDefinition
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, &id); err != nil {
		return err
	}

	for _, idx := range id.Indexes {
		length := len(idx.Columns)
		columns := make([]indexes.Col, length)
		for i, col := range idx.Columns {
			dir := indexes.Ascending
			col.Direction = strings.TrimSpace(strings.ToLower(col.Direction))
			if col.Direction == "desc" || col.Direction == "descending" {
				dir = indexes.Descending
			}
			columns[i] = indexes.Col{
				Name:      col.Name,
				Direction: dir,
			}
		}

		index := indexes.Index{
			Name:    idx.Name,
			Type:    parseIndexType(idx.Type),
			Columns: columns,
		}

		if exists, err := isIndexExists(
			ctx,
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
		if err := iv.CreateOne(ctx, index); err != nil {
			return err
		}
	}
	return nil
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
