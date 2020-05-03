package sqlike

import (
	"context"

	"github.com/Masterminds/semver"
	sqldialect "github.com/si3nloong/sqlike/sql/dialect"
	sqldriver "github.com/si3nloong/sqlike/sql/driver"
	sqlstmt "github.com/si3nloong/sqlike/sql/stmt"
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
func (idv *IndexView) List(ctx context.Context) ([]Index, error) {
	return idv.tb.ListIndexes(ctx)
}

// CreateOne :
func (idv *IndexView) CreateOne(ctx context.Context, idx indexes.Index) error {
	return idv.Create(ctx, []indexes.Index{idx})
}

// Create :
func (idv *IndexView) Create(ctx context.Context, idxs []indexes.Index) error {
	for _, idx := range idxs {
		if len(idx.Columns) < 1 {
			return ErrNoColumn
		}
	}
	stmt := sqlstmt.AcquireStmt(idv.tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	idv.tb.dialect.CreateIndexes(stmt, idv.tb.dbName, idv.tb.name, idxs, idv.isSupportDesc())
	_, err := sqldriver.Execute(
		ctx,
		idv.tb.driver,
		stmt,
		idv.tb.logger,
	)
	return err
}

// CreateOneIfNotExists :
func (idv *IndexView) CreateOneIfNotExists(ctx context.Context, idx indexes.Index) error {
	return idv.CreateIfNotExists(ctx, []indexes.Index{idx})
}

// CreateIfNotExists :
func (idv *IndexView) CreateIfNotExists(ctx context.Context, idxs []indexes.Index) error {
	cols := make([]indexes.Index, 0, len(idxs))
	stmt := sqlstmt.AcquireStmt(idv.tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	for _, idx := range idxs {
		if len(idx.Columns) < 1 {
			return ErrNoColumn
		}
		idv.tb.dialect.HasIndex(stmt, idv.tb.dbName, idv.tb.name, idx)
		var count int
		if err := sqldriver.QueryRowContext(
			ctx,
			idv.tb.driver,
			stmt,
			idv.tb.logger,
		).Scan(&count); err != nil {
			return err
		}
		stmt.Reset()
		if count > 0 {
			continue
		}
		cols = append(cols, idx)
	}
	if len(cols) < 1 {
		return nil
	}
	idv.tb.dialect.CreateIndexes(stmt, idv.tb.dbName, idv.tb.name, cols, idv.isSupportDesc())
	_, err := sqldriver.Execute(
		ctx,
		idv.tb.driver,
		stmt,
		idv.tb.logger,
	)
	return err
}

// DropOne :
func (idv IndexView) DropOne(ctx context.Context, name string) error {
	stmt := sqlstmt.AcquireStmt(idv.tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	idv.tb.dialect.DropIndexes(stmt, idv.tb.dbName, idv.tb.name, []string{name})
	_, err := sqldriver.Execute(
		ctx,
		idv.tb.driver,
		stmt,
		idv.tb.logger,
	)
	return err
}

// DropAll :
func (idv *IndexView) DropAll(ctx context.Context) error {
	idxs, err := idv.List(ctx)
	if err != nil {
		return err
	}
	names := make([]string, 0)
	for _, idx := range idxs {
		names = append(names, idx.Name)
	}
	stmt := sqlstmt.AcquireStmt(idv.tb.dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	idv.tb.dialect.DropIndexes(stmt, idv.tb.dbName, idv.tb.name, names)
	if _, err := sqldriver.Execute(
		ctx,
		idv.tb.driver,
		stmt,
		idv.tb.logger,
	); err != nil {
		return err
	}
	return nil
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

func isIndexExists(ctx context.Context, dbName, table, indexName string, driver sqldriver.Driver, dialect sqldialect.Dialect, logger logs.Logger) (bool, error) {
	stmt := sqlstmt.AcquireStmt(dialect)
	defer sqlstmt.ReleaseStmt(stmt)
	dialect.HasIndexByName(stmt, dbName, table, indexName)
	var count int
	if err := sqldriver.QueryRowContext(
		ctx,
		driver,
		stmt,
		logger,
	).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
