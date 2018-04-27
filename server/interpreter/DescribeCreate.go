package interpreter

import "strings"

// DescribeCreate ...
///////////////////////////////////////////////////////////////////////////////////////////////
func DescribeCreate(tokens []string) *CreateOptions {
	tq := newQueue(tokens).Next()

	cType := strings.ToLower(tq.Current())
	tq.Next()

	ifNotExists := getIfNotExists(tq)

	tableRef := tq.Current()
	tq.Next().Next()

	columnDefs := []*ColumnDef{}
	for {
		columnDefs = append(columnDefs, getColumnDef(tq))
		if tq.Current() == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return &CreateOptions{
		Type:        cType,
		IfNotExists: ifNotExists,
		TableRef:    tableRef,
		ColumnDefs:  columnDefs,
	}
}
