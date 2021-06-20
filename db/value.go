package db

import "context"

// SqlValuer :
type SqlValuer interface {
	SqlValue(ctx context.Context, v interface{}) error
}
