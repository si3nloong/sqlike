package sqlike

import (
	"context"

	"github.com/blang/semver"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/sqlike/indexes"
	"github.com/si3nloong/sqlike/sqlike/logs"
	"github.com/si3nloong/sqlike/types"
	"errors"
)

// Index :
type Index struct {
	Name      string
	Type      string
	IsVisible types.Boolean
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
	return idv.CreateMany([]indexes.Index{idx})
}

// CreateMany :
func (idv *IndexView) CreateMany(idxs []indexes.Index) error {
	for _, idx := range idxs {
		if len(idx.Columns) < 1 {
			return errors.New("sqlike: empty columns to create index")
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
	mysql8 := semver.MustParse("8.0.0")
	version := idv.tb.client.Version()
	flag := false
	if idv.tb.client.driverName == "mysql" && version.GTE(mysql8) {
		flag = true
	}
	idv.supportDesc = &flag
	return *idv.supportDesc
}

func isIndexExists(dbName, table, indexName string, driver sqldriver.Driver, dialect sqldialect.Dialect, logger logs.Logger) bool {
	var count int
	if err := sqldriver.QueryRowContext(
		context.Background(),
		driver,
		dialect.HasIndex(dbName, table, indexName),
		logger,
	).Scan(&count); err != nil {
		panic(err)
	}
	return count > 0
}
