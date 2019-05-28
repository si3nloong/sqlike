package sqlike

import "database/sql"

// Session :
type Session interface {
}

// Transaction :
type Transaction struct {
	tx *sql.Tx
}

// RollbackTransaction :
func (tx *Transaction) RollbackTransaction() error {
	return tx.tx.Rollback()
}

// CommitTransaction :
func (tx *Transaction) CommitTransaction() error {
	return tx.tx.Commit()
}
