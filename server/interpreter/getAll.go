package interpreter

import "strings"

func getAll(tq *queue) bool {
	if strings.ToLower(tq.Current()) == ALL {
		tq.Next()
		return true
	}
	return false
}
