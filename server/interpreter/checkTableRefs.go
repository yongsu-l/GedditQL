package interpreter

import "strings"

func checkTableRefs(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	for {
		err := checkName(tq)
		if err != nil {
			return newError(ErrIllTabRef, tq)
		}

		if strings.ToLower(tq.Current()) == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return nil
}
