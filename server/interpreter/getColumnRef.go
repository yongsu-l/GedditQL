package interpreter

import "strings"

func getColumnRef(tq *queue) string {
	s := tq.Current()
	if strings.ToLower(tq.Next().Current()) == DOT {
		s += tq.Current() + tq.Next().Current()
		tq.Next()
	}
	return s
}
