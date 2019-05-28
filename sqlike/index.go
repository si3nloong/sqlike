package sqlike

import (
	"strings"

	"github.com/si3nloong/sqlike/sqlike/indexes"
	sqldriver "github.com/si3nloong/sqlike/sqlike/sql/driver"
	"github.com/si3nloong/sqlike/types"
	"golang.org/x/xerrors"
)

// Index :
type Index struct {
	Name      string
	Type      string
	IsVisible types.Boolean
}

// IndexView :
type IndexView struct {
	tb     *Table
	logger Logger
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
	for i, idx := range idxs {
		if len(idx.Columns) < 1 {
			return xerrors.New("empty columns to create index")
		}
		if idx.Name == "" {
			idxs[i].Name = strings.Join(idx.Columns, "_") + "_idx"
		}
	}
	_, err := sqldriver.Execute(
		idv.tb.driver,
		idv.tb.dialect.CreateIndexes(idv.tb.name, idxs),
		idv.logger,
	)
	return err
}

// DropOne :
func (idv IndexView) DropOne(name string) error {
	_, err := sqldriver.Execute(
		idv.tb.driver,
		idv.tb.dialect.DropIndex(idv.tb.name, name),
		idv.logger,
	)
	return err
}
