package interpreter

import "strings"

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
