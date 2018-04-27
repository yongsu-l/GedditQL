package interpreter

import "strings"

func checkSetExprs(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrNoSetExp, tq)
	}

	for {
		err := checkColumnRef(tq)
		if err != nil {
			return err
		}

		if strings.ToLower(tq.Current()) != EQ {
			return newError(ErrUnexpToken, tq)
		}

		b, _, err := checkValue(tq.Next())

		if err != nil {
			return err
		}

		if !b {
			return newError(ErrIllValue, tq)
		}

		if strings.ToLower(tq.Current()) == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return nil
}
