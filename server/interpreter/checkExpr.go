package interpreter

import "strings"

func checkExpr(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrIllEOQuery, tq)
	}
	err := checkTerm(tq)
	if err != nil {
		return err
	}

	switch tk := strings.ToLower(tq.Current()); tk {
	case LESS, ISNOT, LESSEQ, GREATER, GREATEREQ, EQ, NEQ:
		err = checkTerm(tq.Next())
		if err != nil {
			return err
		}
		break
	default:
		return newError(ErrUnexpToken, tq)
	}
	return nil
}
