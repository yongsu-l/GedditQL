package interpreter

import "strings"

func checkColumnRef(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	err := checkName(tq)
	if err != nil {
		return newError(ErrIllColRef, tq)
	}

	if strings.ToLower(tq.Current()) == DOT {
		if strings.ToLower(tq.Next().Current()) == EOQ {
			return newError(ErrIllEOQuery, tq)
		}

		err = checkName(tq)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
