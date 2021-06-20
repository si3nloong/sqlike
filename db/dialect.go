package db

// SQLDialect :
type SQLDialect interface {
	TableName(db, table string) string
	Var(i int) string
	Quote(n string) string
	Format(v interface{}) (val string)
}
