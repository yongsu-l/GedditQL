package interpreter

import (
	"bytes"
	"strconv"
	"strings"
)

// var
var (
	buffer    bytes.Buffer
	PopBuffer = func() string {
		s := buffer.String()
		buffer.Reset()
		return s
	}
)

// DescribeSelect inputs a list of tokens of a select statement and returns extracted parameters
///////////////////////////////////////////////////////////////////////////////////////////////
func DescribeSelect(tokens []string) *SelectOptions {
	// puts tokens into a queue
	tq := newQueue(tokens)
	// skips the SELECT token
	tq.Next()

	// check for DISTINCT
	distinct := getDistinct(tq)
	// check for *
	all := getAll(tq)

	// get the column refs
	columnRefs := []string{}
	as := map[string]string{}
	sa := map[string]string{}
	if !all {
		columnRefs, as, sa = getColumnRefs(tq)
	}

	// get the table refs
	tableRefs := []string{}
	if strings.ToLower(tq.Current()) == FROM {
		tableRefs = getTableRefs(tq.Next())
	}

	// get the WHERE clause
	var condition func(map[string]string) (bool, error)
	if strings.ToLower(tq.Current()) == WHERE {
		condition = getCondition(tq.Next(), sa)
	}

	// get the ORDER BY clause
	order, by := "", ASC
	if strings.ToLower(tq.Current()) == ORDER {
		order = strings.ToLower(tq.Next().Next().Current())
		tq.Next()

		// checks for the optional ASC or DESC (default to ASC if absent)
		if strings.ToLower(tq.Current()) == ASC || strings.ToLower(tq.Current()) == DESC {
			by = strings.ToLower(tq.Current())
			tq.Next()
		}
	}

	// get the LIMIT clause
	limit := 0
	if strings.ToLower(tq.Current()) == LIMIT {
		limit, _ = strconv.Atoi(strings.ToLower(tq.Next().Current()))
	}

	return &SelectOptions{
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
	if strings.ToLower(tq.Current()) == DISTINCT {
		tq.Next()
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getAll(tq *queue) bool {
	if strings.ToLower(tq.Current()) == ALL {
		tq.Next()
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getColumnRef(tq *queue) string {
	s := tq.Current()
	if strings.ToLower(tq.Next().Current()) == DOT {
		s += tq.Current() + tq.Next().Current()
		tq.Next()
	}
	return s
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getColumnRefs(tq *queue) ([]string, map[string]string, map[string]string) {
	columnRefs := []string{}
	// columnRef: AS
	as := map[string]string{}
	// AS: columnRef
	sa := map[string]string{}

	for {
		s := getColumnRef(tq)
		columnRefs = append(columnRefs, s)
		// checks for AS clause
		if strings.ToLower(tq.Current()) == AS {
			a := tq.Next().Current()
			as[s] = a
			sa[a] = s
			tq.Next()
		}

		// checks if there are more columnRefs
		if strings.ToLower(tq.Current()) == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return columnRefs, as, sa
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getTableRefs(tq *queue) []string {
	tableRefs := []string{}

	for {
		tableRefs = append(tableRefs, tq.Current())
		if strings.ToLower(tq.Next().Current()) == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return tableRefs
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getCondition(tq *queue, sa map[string]string) func(map[string]string) (bool, error) {
	res := getExpr(tq, sa)

	for strings.ToLower(tq.Current()) == AND || strings.ToLower(tq.Current()) == OR {
		op := strings.ToLower(tq.Current())
		next := getExpr(tq.Next(), sa)
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
func getExpr(tq *queue, sa map[string]string) func(map[string]string) (bool, error) {
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

		val, ok := sa[val1]
		if ok {
			val1 = val
		}
	}

	op := strings.ToLower(tq.Current())

	isVal2, val2 := getValue(tq.Next())
	if !isVal2 {
		val2 = getColumnRef(tq)
		curr2 = val2
		pos2 = tq.Pos()

		val, ok := sa[val2]
		if ok {
			val2 = val
		}
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
	if strings.ToLower(tq.Current()) == SUB {
		isNumeric, val := getNumeric(tq.Next())
		return isNumeric, SUB + val
	}

	isNumeric, val := getNumeric(tq)
	if isNumeric {
		return isNumeric, val
	}

	if strings.ToLower(tq.Current()) == TRUE || strings.ToLower(tq.Current()) == FALSE {
		val = strings.ToLower(tq.Current())
		tq.Next()
		return true, val
	}

	if strings.ToLower(tq.Current())[0] == '\'' || strings.ToLower(tq.Current())[0] == '"' {
		val = strings.ToLower(tq.Current())
		tq.Next()
		return true, val
	}

	return false, ""
}

///////////////////////////////////////////////////////////////////////////////////////////////
func getNumeric(tq *queue) (bool, string) {
	if _, err := strconv.Atoi(strings.ToLower(tq.Current())); err == nil {
		n := strings.ToLower(tq.Current())
		if strings.ToLower(tq.Next().Current()) == DOT {
			n = n + DOT + strings.ToLower(tq.Next().Current())
			tq.Next()
			return true, n
		}

		return true, n
	}

	return false, ""
}
