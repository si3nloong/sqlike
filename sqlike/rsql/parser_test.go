package rsql

import "testing"

type testStruct struct {
	ID   string
	Name string
}

func TestParser(t *testing.T) {
	p := MustNewParser(testStruct{})
	p.ParseQuery([]byte(`(_id==133,(category!=-10.00;num==.922;test=="value\""));d1=="";c1==testing,d1!=108)`))
	// 	p := new(Parser)
	// 	p.SetComparisonOperator(BasicOperator([]string{">", "=gt="}, false))
	// 	p.SetComparisonOperator(BasicOperator([]string{"<", "=lt="}, false))
	// 	p.SetComparisonOperator(BasicOperator([]string{">=", "=ge="}, false))
	// 	p.SetComparisonOperator(BasicOperator([]string{"<=", "=le="}, false))
	// 	p.SetComparisonOperator(BasicOperator([]string{"==", "=eq="}, false))
	// 	p.SetComparisonOperator(BasicOperator([]string{"!=", "=ne="}, false))
	// 	p.SetComparisonOperator(BasicOperator([]string{"=in="}, true))
	// 	p.SetComparisonOperator(BasicOperator([]string{"=nin="}, true))
}
