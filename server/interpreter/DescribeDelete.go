package interpreter

import (
	"strconv"
	"strings"
)

// DescribeDelete ...
///////////////////////////////////////////////////////////////////////////////////////////////
func DescribeDelete(tokens []string) *DeleteOptions {
	tq := newQueue(tokens).Next().Next()
	tableRef := tq.Current()

	var condition func(map[string]string) (bool, error)
	if strings.ToLower(tq.Next().Current()) == WHERE {
		condition = getCondition(tq.Next(), map[string]string{})
	}

	var limit int
	if strings.ToLower(tq.Current()) == LIMIT {
		limit, _ = strconv.Atoi(strings.ToLower(tq.Next().Current()))
	}

	return &DeleteOptions{
		TableRef:  tableRef,
		Condition: condition,
		Limit:     limit,
	}
}
