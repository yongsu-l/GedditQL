package interpreter

import (
	"strconv"
	"strings"
)

// DescribeSelect ...
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
	funcCols := []string{}
	funcMap := map[string]string{}

	if !all {
		columnRefs, as, sa, funcCols, funcMap = getSelectExprs(tq)
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
		FuncCols:   funcCols,
		FuncMap:    funcMap,
		TableRefs:  tableRefs,
		Condition:  condition,
		Order:      order,
		By:         by,
		Limit:      limit,
	}
}
