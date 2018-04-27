package interpreter

import (
	"bytes"
	"strconv"
)

// var
///////////////////////////////////////////////////////////////////////////////////////////////
var (
	buffer    bytes.Buffer
	PopBuffer = func() string {
		s := buffer.String()
		buffer.Reset()
		return s
	}
)

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
func describeSelect(tokens []string) *selectOptions {
	tq := newQueue(tokens)
	tq.Next()

	distinct := getDistinct(tq)
	all := getAll(tq)

	columnRefs := []string{}
	as := map[string]string{}
	if !all {
		columnRefs, as = getColumnRefs(tq)
	}

	tableRefs := []string{}
	if tq.Current() == FROM {
		tableRefs = getTableRefs(tq.Next())
	}

	var condition func(map[string]string) (bool, error)
	if tq.Current() == WHERE {
		condition = getCondition(tq.Next())
	}

	order, by := "", ASC
	if tq.Current() == ORDER {
		order = tq.Next().Next().Current()
		tq.Next()
		if tq.Current() == ASC || tq.Current() == DESC {
			by = tq.Current()
			tq.Next()
		}
	}

	limit := -1
	if tq.Current() == LIMIT {
		limit, _ = strconv.Atoi(tq.Next().Current())
	}

	return &selectOptions{
		Distinct:   distinct,
		All:        all,
		ColumnRefs: columnRefs,
		As:         as,
		TableRefs:  tableRefs,
		Condition:  condition,
		Order:      order,
		By:         by,
		Limit:      limit,
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getDistinct(tq *queue) bool {
	if tq.Current() == DISTINCT {
		tq.Next()
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getAll(tq *queue) bool {
	if tq.Current() == ALL {
		tq.Next()
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getColumnRef(tq *queue) string {
	s := tq.Current()
	if tq.Next().Current() == DOT {
		s += tq.Current() + tq.Next().Current()
		tq.Next()
	}
	return s
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getColumnRefs(tq *queue) ([]string, map[string]string) {
	columnRefs := []string{}
	as := map[string]string{}

	for {
		s := getColumnRef(tq)
		columnRefs = append(columnRefs, s)
		if tq.Current() == AS {
			a := tq.Next().Current()
			as[s] = a
			as[a] = s
			tq.Next()
		} else {
			as[s] = s
		}

		if tq.Current() == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return columnRefs, as
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getTableRefs(tq *queue) []string {
	tableRefs := []string{}

	for {
		tableRefs = append(tableRefs, tq.Current())
		if tq.Next().Current() == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return tableRefs
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getCondition(tq *queue) func(map[string]string) (bool, error) {
	res := getExpr(tq)

	for tq.Current() == AND || tq.Current() == OR {
		op := tq.Current()
		next := getExpr(tq.Next())
		prev := res

		res = func(dict map[string]string) (bool, error) {
			b1, err := prev(dict)
			if err != nil {
				return false, err
			}
			b2, err := next(dict)
			if err != nil {
				return false, err
			}

			switch op {
			case AND:
				return b1 && b2, nil
			case OR:
				return b1 || b2, nil
			default:
				return false, nil
			}
		}
	}

	return res
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getExpr(tq *queue) func(map[string]string) (bool, error) {
	var (
		curr1 string
		pos1  int
		curr2 string
		pos2  int
	)

	isVal1, val1 := getValue(tq)
	if !isVal1 {
		val1 = getColumnRef(tq)
		curr1 = val1
		pos1 = tq.Pos()
	}

	op := tq.Current()

	isVal2, val2 := getValue(tq.Next())
	if !isVal2 {
		val2 = getColumnRef(tq)
		curr2 = val2
		pos2 = tq.Pos()
	}

	return func(dict map[string]string) (bool, error) {
		if !isVal1 {
			val, ok := dict[val1]
			if !ok {
				return false, newError2(ErrNoColRef, curr1, pos1)
			}
			val1 = val
		}

		if !isVal2 {
			val, ok := dict[val2]
			if !ok {
				return false, newError2(ErrNoColRef, curr2, pos2)
			}
			val2 = val
		}

		switch op {
		case "<":
			return val1 < val2, nil
		case "<>", "!=":
			return val1 != val2, nil
		case "<=":
			return val1 <= val2, nil
		case ">":
			return val1 > val2, nil
		case ">=":
			return val1 >= val2, nil
		case "=":
			return val1 == val2, nil
		default:
			return false, nil
		}
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getValue(tq *queue) (bool, string) {
	if tq.Current() == SUB {
		isNumeric, val := getNumeric(tq.Next())
		return isNumeric, SUB + val
	}

	isNumeric, val := getNumeric(tq)
	if isNumeric {
		return isNumeric, val
	}

	if tq.Current() == TRUE || tq.Current() == FALSE {
		val = tq.Current()
		tq.Next()
		return true, val
	}

	if tq.Current()[0] == '\'' || tq.Current()[0] == '"' {
		val = tq.Current()
		tq.Next()
		return true, val
	}

	return false, ""
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getNumeric(tq *queue) (bool, string) {
	if _, err := strconv.Atoi(tq.Current()); err == nil {
		n := tq.Current()
		if tq.Next().Current() == DOT {
			n = n + DOT + tq.Next().Current()
			tq.Next()
			return true, n
		}

		return true, n
	}

	return false, ""
}

// func Distinct(tokens []string) (bool, []string, error) {
// 	switch tq.Show() {
// 	case EOF:
// 		return false, nil, errors.New("Illegal End of Query")
// 	case DISTINCT:
// 		return true, tokens[1:], nil
// 	default:
// 		return false, tokens, nil
// 	}
// }

// func All(tokens []string) (bool, []string, error) {
// 	switch tq.Show() {
// 	case EOF:
// 		return false, nil, errors.New("Illegal End of Query")
// 	case ALL:
// 		return true, tokens[1:], nil
// 	default:
// 		return false, tokens, nil
// 	}
// }

// func Value(tokens []string) (*class.Value, []string, error) {
// 	switch tq.Show() {
// 	case EOF:
// 		return nil, nil, errors.New("Unexpected End of Query")
// 	case TRUE, FALSE:
// 		return &class.Value{Type: "bool", Value: tq.Show()}, tokens[1:], nil
// 	default:
// 		if _, err := strconv.Atoi(tq.Show()); err == nil {
// 			buffer.WriteString(tq.Show())
// 			tokens = tokens[1:]
// 			if tq.Show() == "." {
// 				buffer.WriteString(tq.Show())
// 				if _, err := strconv.Atoi(tokens[1]); err == nil {
// 					tokens = tokens[1:]
// 					buffer.WriteString(tq.Show())
// 					tokens = tokens[1:]
// 					return &class.Value{Type: "float", Value: PopBuffer()}, tokens, nil
// 				} else {
// 					return nil, nil, errors.New("Unexpected Token: " + tq.Show())
// 				}
// 			} else {
// 				return &class.Value{Type: "int", Value: PopBuffer()}, tokens, nil
// 			}
// 		}

// 		if tq.Show()[0] == '\'' && tq.Show()[len(tq.Show())-1] == '\'' {
// 			return &class.Value{Type: "string", Value: tq.Show()}, tokens[1:], nil
// 		}

// 		if tq.Show()[0] == '"' && tq.Show()[len(tq.Show())-1] == '"' {
// 			return &class.Value{Type: "string", Value: tq.Show()}, tokens[1:], nil
// 		}

// 		return nil, tokens, nil
// 	}
// }

// func ColumnRef(tokens []string) (*class.ColumnRef, []string, error) {
// 	switch tq.Show() {
// 	case EOF:
// 		return nil, nil, errors.New("Unexpected End of Query")
// 	default:
// 		cr := tq.Show()
// 		tokens = tokens[1:]

// 		tr := ""
// 		if tq.Show() == "." {
// 			tr = cr
// 			tokens = tokens[1:]
// 			if tq.Show() == EOF {
// 				return nil, nil, errors.New("Unexpected End of Query")
// 			} else {
// 				cr = tq.Show()
// 				tokens = tokens[1:]
// 			}
// 		}

// 		return &class.ColumnRef{TableRef: tr, Name: cr}, tokens, nil
// 	}
// }

// func Term(tokens []string) (*class.Term, []string, error) {
// 	if tq.Show() == EOF {
// 		return nil, nil, errors.New("Unexpected End of Query")
// 	}

// 	v, tokens, err := Value(tokens)

// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	if v != nil {
// 		return &class.Term{Value: v}, tokens, nil
// 	}

// 	cr, tokens, err := ColumnRef(tokens)

// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	if cr != nil {
// 		return &class.Term{ColumnRef: cr}, tokens, nil
// 	}

// 	return nil, nil, errors.New("Expected Value or Column Reference At: " + tq.Show())
// }

// func SelectExprs(tokens []string) ([]*class.SelectExpr, []string, error) {
// 	if tq.Show() == FROM || tq.Show() == EOF {
// 		return make([]*class.SelectExpr, 0), tokens, nil
// 	}

// 	t, tokens, err := Term(tokens)

// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	as := ""
// 	if tq.Show() == AS {
// 		tokens = tokens[1:]
// 		as = tq.Show()
// 		tokens = tokens[1:]
// 	}

// 	switch tq.Show() {
// 	case COMMA, FROM, EOF:
// 		break
// 	default:
// 		return nil, nil, errors.New("Unexpected Token: " + tq.Show())
// 	}

// 	tokens = tokens[1:]

// 	res := make([]*class.SelectExpr, 1)
// 	res[0] = &class.SelectExpr{Term: t, As: as}
// 	rest, tokens, err := SelectExprs(tokens)

// 	return append(res, rest...), tokens, err
// }

// func TableRefs(tokens []string) ([]string, []string, error) {
// 	switch tq.Show() {
// 	case EOF, WHERE, ORDER, LIMIT:
// 		return make([]string, 0), tokens, nil
// 	default:
// 		res := tq.Show()
// 		tokens = tokens[1:]

// 		switch tq.Show() {
// 		case COMMA:
// 			tokens = tokens[1:]
// 		case EOF, WHERE, ORDER, LIMIT:
// 			break
// 		default:
// 			return nil, nil, errors.New("Unexpected Token: " + tq.Show())
// 		}

// 		rest, tokens, err := TableRefs(tokens)
// 		return append([]string{res}, rest...), tokens, err
// 	}
// }

// func Expr(tokens []string) (*class.Expr, []string, error) {
// 	t, tokens, err := Term(tokens)

// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	switch tq.Show() {
// 	case "+", "-", "*", "/", "%":
// 		op := tq.Show()
// 		rExpr, tokens, err := Expr(tokens[1:])
// 		return &class.Expr{Term: t, Op: op, RExpr: rExpr}, tokens, err
// 	default:
// 		return &class.Expr{Term: t, Op: ""}, tokens, err
// 	}
// }

// func Condition(tokens []string) (*class.Condition, []string, error) {
// 	switch tq.Show() {
// 	case EOF:
// 		return nil, nil, errors.New("Illegal End of Query")
// 		case
// 	}
// }
