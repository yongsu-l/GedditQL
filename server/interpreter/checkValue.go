package interpreter

import (
	"strconv"
	"strings"
)

func checkValue(tq *queue) (bool, string, error) {
	if strings.ToLower(tq.Current()) == EOQ {
		return false, "", newError(ErrIllEOQuery, tq)
	}

	if strings.ToLower(tq.Current()) == SUB {
		b, typ, err := checkNumeric(tq.Next())
		if !b || err != nil {
			return false, "", newError(ErrIllValue, tq)
		}
		return true, typ, nil
	}

	b, typ, err := checkNumeric(tq)
	if err != nil {
		return false, "", err
	}
	if b {
		return true, typ, nil
	}

	if strings.ToLower(tq.Current())[0] == '"' {
		if strings.ToLower(tq.Current())[len(strings.ToLower(tq.Current()))-1] == '"' {
			tq.Next()
			return true, "string", nil
		}
		return false, "", newError(ErrOpenQuote, tq)
	}

	if strings.ToLower(tq.Current()) == TRUE || strings.ToLower(tq.Current()) == FALSE {
		tq.Next()
		return true, "bool", nil
	}

	return false, "", nil
}

func checkNumeric(tq *queue) (bool, string, error) {
	if _, err := strconv.Atoi(strings.ToLower(tq.Current())); err == nil {
		if strings.ToLower(tq.Next().Current()) != DOT {
			return true, "int", nil
		}

		if _, err := strconv.Atoi(strings.ToLower(tq.Next().Current())); err == nil {
			tq.Next()
			return true, "float", nil
		}

		return false, "", newError(ErrIllValue, tq)
	}

	return false, "", nil
}
