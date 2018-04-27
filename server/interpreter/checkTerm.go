package interpreter

import (
	"strings"
)

func checkTerm(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	b, _, err := checkValue(tq)

	if err != nil {
		return err
	}
	if b {
		return nil
	}

	return checkColumnRef(tq)
}
