package rsql

import (
	"github.com/si3nloong/sqlike/sql/expr"
	"github.com/si3nloong/sqlike/sqlike/primitive"
)

func (p *Parser) parseFilter(values map[string]string, params *Params) (errs Errors) {
	val, ok := values[p.FilterTag]
	if !ok || len(val) < 1 {
		return
	}

	lxr, _ := p.lexer.Scanner([]byte(val))
	scan := &Scanner{parser: p, Scanner: lxr}
	grp := new(primitive.Group)
	stack := make([]bool, 0)

	for {
		tkn, eof := scan.NextToken()
		if eof {
			break
		}

		char := string(tkn.Lexeme)
		switch tkn.Type {
		case Whitespace:
			// skip
		case Group:
			{
				if char == "(" {
					stack = append(stack, true)
					grp.Values = append(grp.Values, expr.Raw("("))
					continue
				}

				if len(stack) < 0 {
					errs = append(errs, &FieldError{Module: p.FilterTag})
					continue
				}

				stack = stack[:len(stack)-1]
				grp.Values = append(grp.Values, expr.Raw(")"))
			}

		case Text:
			// check expression
			if err := scan.ParseExpression(grp, tkn); err != nil {

			}

		case String:
			if err := scan.ParseExpression(grp, tkn); err != nil {

			}
		case And:
			grp.Values = append(grp.Values, primitive.And)
		case Or:
			grp.Values = append(grp.Values, primitive.Or)
		}
	}

	params.Filters = *grp
	return
}
