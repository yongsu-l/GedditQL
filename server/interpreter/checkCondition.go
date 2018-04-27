package interpreter

import "strings"

func checkCondition(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	err := checkExpr(tq)
	if err != nil {
		return err
	}

	if strings.ToLower(tq.Current()) == AND || strings.ToLower(tq.Current()) == OR {
		err := checkCondition(tq.Next())
		if err != nil {
			return err
		}
	}

	return nil
}
