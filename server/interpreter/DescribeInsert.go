package interpreter

import "strings"

// DescribeInsert ...
///////////////////////////////////////////////////////////////////////////////////////////////
func DescribeInsert(tokens []string) *InsertOptions {
	tq := newQueue(tokens).Next().Next()

	tableRef := tq.Current()

	columnRefs := []string{}
	if strings.ToLower(tq.Next().Current()) != VALUES {
		tq.Next()
		for {
			columnRefs = append(columnRefs, tq.Current())

			if tq.Next().Current() == COMMA {
				tq.Next()
			} else {
				break
			}
		}
		tq.Next()
	}

	tq.Next().Next()
	values := []string{}
	types := []string{}
	for {
		_, val, typ := getValue(tq)
		values = append(values, val)
		types = append(types, typ)

		if tq.Current() == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return &InsertOptions{
		TableRef:   tableRef,
		ColumnRefs: columnRefs,
		Values:     values,
		ValueTypes: types,
	}
}
