package sqlike

import (
	"context"

	"github.com/Masterminds/semver"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/si3nloong/sqlike/sqlike/logs"
)

var mysql8 = semver.MustParse("8.0.0")

// Index :
type Index struct {
	Name     string
	Type     string
	IsUnique bool
}

// IndexView :
type IndexView struct {
	tb          *Table
	supportDesc *bool
}

// List :
func (idv *IndexView) List() ([]Index, error) {
	return idv.tb.ListIndexes()
}

// CreateOne :
func (idv *IndexView) CreateOne(idx indexes.Index) error {
	return idv.Create([]indexes.Index{idx})
}

// Create :
func (idv *IndexView) Create(idxs []indexes.Index) error {
	for _, idx := range idxs {
		if len(idx.Columns) < 1 {
			return ErrNoColumn
		}
	}
	_, err := sqldriver.Execute(
		context.Background(),
		idv.tb.driver,
		idv.tb.dialect.CreateIndexes(idv.tb.dbName, idv.tb.name, idxs, idv.isSupportDesc()),
		idv.tb.logger,
	)
	return err
}

// CreateOneIfNotExists :
func (idv *IndexView) CreateOneIfNotExists(idx indexes.Index) error {
	return idv.CreateIfNotExists([]indexes.Index{idx})
}

// CreateIfNotExists :
func (idv *IndexView) CreateIfNotExists(idxs []indexes.Index) error {
	cols := make([]indexes.Index, 0, len(idxs))
	for _, idx := range idxs {
		if len(idx.Columns) < 1 {
			return ErrNoColumn
		}
		var count int
		if err := sqldriver.QueryRowContext(
			context.Background(),
			idv.tb.driver,
			idv.tb.dialect.HasIndex(idv.tb.dbName, idv.tb.name, idx),
			idv.tb.logger,
		).Scan(&count); err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		cols = append(cols, idx)
	}
	if len(cols) < 1 {
		return nil
	}
	_, err := sqldriver.Execute(
		context.Background(),
		idv.tb.driver,
		idv.tb.dialect.CreateIndexes(idv.tb.dbName, idv.tb.name, cols, idv.isSupportDesc()),
		idv.tb.logger,
	)
	return err
}

// DropOne :
func (idv IndexView) DropOne(name string) error {
	_, err := sqldriver.Execute(
		context.Background(),
		idv.tb.driver,
		idv.tb.dialect.DropIndex(idv.tb.dbName, idv.tb.name, name),
		idv.tb.logger,
	)
	return err
}

func (idv *IndexView) isSupportDesc() bool {
	if idv.supportDesc != nil {
		return *idv.supportDesc
	}
	flag := false
	if idv.tb.client.driverName == "mysql" &&
		idv.tb.client.version.GreaterThan(mysql8) {
		flag = true
	}
	idv.supportDesc = &flag
	return *idv.supportDesc
}

func isIndexExists(dbName, table, indexName string, driver sqldriver.Driver, dialect sqldialect.Dialect, logger logs.Logger) (bool, error) {
	var count int
	if err := sqldriver.QueryRowContext(
		context.Background(),
		driver,
		dialect.HasIndexByName(dbName, table, indexName),
		logger,
	).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
