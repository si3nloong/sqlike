package expr

import (
	"github.com/si3nloong/sqlike/sql"

	"github.com/si3nloong/sqlike/sqlike/primitive"
)

// Union :
func Union(stmt1 *sql.SelectStmt, stmt2 *sql.SelectStmt, others ...*sql.SelectStmt) (grp primitive.Group) {
	grp = union(Raw("UNION"), append([]*sql.SelectStmt{stmt1, stmt2}, others...))
	return
}

// UnionAll :
func UnionAll(stmt1 *sql.SelectStmt, stmt2 *sql.SelectStmt, others ...*sql.SelectStmt) (grp primitive.Group) {
	grp = union(Raw("UNION ALL"), append([]*sql.SelectStmt{stmt1, stmt2}, others...))
	return
}

func union(link primitive.Raw, stmts []*sql.SelectStmt) (grp primitive.Group) {
	for i, stmt := range stmts {
		if i > 0 {
			grp.Values = append(grp.Values, link)
		}
		grp.Values = append(grp.Values, stmt)
	}
	return
}
