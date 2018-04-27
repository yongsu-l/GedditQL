package interpreter

import (
	"strings"
)

// DescribeUpdate ...
///////////////////////////////////////////////////////////////////////////////////////////////
func DescribeUpdate(tokens []string) *UpdateOptions {
	tq := newQueue(tokens).Next()

	tableRef := tq.Current()
	valueMap, typeMap := getSetExprs(tq.Next().Next())

	var condition func(map[string]string) (bool, error)
	if strings.ToLower(tq.Current()) == WHERE {
		condition = getCondition(tq.Next(), map[string]string{})
	}

	return &UpdateOptions{
		TableRef:  tableRef,
		ValueMap:  valueMap,
		TypeMap:   typeMap,
		Condition: condition,
	}
}
