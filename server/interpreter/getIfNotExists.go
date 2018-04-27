package interpreter

import "strings"

func getIfNotExists(tq *queue) bool {
	if strings.ToLower(tq.Current()) == IF {
		tq.Next().Next().Next()
		return true
	}
	return false
}
