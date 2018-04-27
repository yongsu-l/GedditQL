package interpreter

import "strings"

func getDistinct(tq *queue) bool {
	if strings.ToLower(tq.Current()) == DISTINCT {
		tq.Next()
		return true
	}
	return false
}
