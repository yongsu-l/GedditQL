package interpreter

import (
	"strconv"
	"strings"
)

func getValue(tq *queue) (bool, string, string) {
	if strings.ToLower(tq.Current()) == SUB {
		isNumeric, val, typ := getNumeric(tq.Next())
		return isNumeric, SUB + val, typ
	}

	isNumeric, val, typ := getNumeric(tq)
	if isNumeric {
		return isNumeric, val, typ
	}

	if strings.ToLower(tq.Current()) == TRUE || strings.ToLower(tq.Current()) == FALSE {
		val = strings.ToLower(tq.Current())
		tq.Next()
		return true, val, "bool"
	}

	if strings.ToLower(tq.Current())[0] == '"' {
		val = strings.ToLower(tq.Current())
		tq.Next()
		return true, val, "string"
	}

	return false, "", ""
}

func getNumeric(tq *queue) (bool, string, string) {
	if _, err := strconv.Atoi(strings.ToLower(tq.Current())); err == nil {
		n := strings.ToLower(tq.Current())
		if strings.ToLower(tq.Next().Current()) == DOT {
			n = n + DOT + strings.ToLower(tq.Next().Current())
			tq.Next()
			return true, n, "float"
		}

		return true, n, "int"
	}

	return false, "", ""
}
