package interpreter

import "strings"

func checkUpdate(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	err := checkName(tq)
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	if strings.ToLower(tq.Current()) != SET {
		return newError(ErrExpToken+"\"SET\"", tq)
	}

	err = checkSetExprs(tq.Next())
	if err != nil {
		return err
	}

	if strings.ToLower(tq.Current()) == WHERE {
		err := checkCondition(tq.Next())
		if err != nil {
			return err
		}
	}

	if strings.ToLower(tq.Current()) != EOQ {
		return newError(ErrUnexpToken, tq)
	}

	return nil
}
