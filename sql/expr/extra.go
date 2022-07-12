package expr

import "github.com/si3nloong/sqlike/v2/internal/primitive"

type selectStmt interface {
	// Where(fields ...any) selectStmt
}

// Union :
func Union(stmt1 selectStmt, stmt2 selectStmt, others ...selectStmt) (grp primitive.Group) {
	grp = union(Raw(" UNION "), append([]selectStmt{stmt1, stmt2}, others...))
	return
}

// // UnionAll :
// func UnionAll(stmt1 *sql.SelectStmt, stmt2 *sql.SelectStmt, others ...*sql.SelectStmt) (grp primitive.Group) {
// 	grp = union(Raw("UNION ALL"), append([]*sql.SelectStmt{stmt1, stmt2}, others...))
// 	return
// }

// Exists :
func Exists(subquery any) (grp primitive.Group) {
	grp.Values = append(grp.Values, Raw("EXISTS ("))
	grp.Values = append(grp.Values, subquery)
	grp.Values = append(grp.Values, Raw(")"))
	return
}

// NotExists :
func NotExists(subquery any) (grp primitive.Group) {
	grp.Values = append(grp.Values, Raw("NOT EXISTS ("))
	grp.Values = append(grp.Values, subquery)
	grp.Values = append(grp.Values, Raw(")"))
	return
}

// Case :
func Case() *primitive.Case {
	return new(primitive.Case)
}

func union(link primitive.Raw, stmts []selectStmt) (grp primitive.Group) {
	for i, stmt := range stmts {
		if i > 0 {
			grp.Values = append(grp.Values, link)
		}
		grp.Values = append(grp.Values, Raw("("))
		grp.Values = append(grp.Values, stmt)
		grp.Values = append(grp.Values, Raw(")"))
	}
	return
}
