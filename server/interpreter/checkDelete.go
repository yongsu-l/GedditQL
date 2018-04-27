package interpreter

import (
	"strconv"
	"strings"
)

func checkDelete(tq *queue) error {
	if strings.ToLower(tq.Current()) != FROM {
		return newError(ErrUnexpToken, tq)
	}

	err := checkName(tq.Next())
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	if strings.ToLower(tq.Current()) == WHERE {
		err = checkCondition(tq.Next())
		if err != nil {
			return err
		}
	}

	if strings.ToLower(tq.Current()) == LIMIT {
		if _, err := strconv.Atoi(strings.ToLower(tq.Next().Current())); err != nil {
			return newError(ErrExpInt, tq)
		}
		tq.Next()
	}

	if strings.ToLower(tq.Current()) != EOQ {
		return newError(ErrUnexpToken, tq)
	}

	return nil
}
