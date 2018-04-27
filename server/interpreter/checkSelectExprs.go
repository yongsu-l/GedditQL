package interpreter

import "strings"

func checkSelectExprs(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrIllEOQuery, tq)
	}
	if strings.ToLower(tq.Current()) == FROM {
		return newError(ErrNoSelExp, tq)
	}

	for {
		err := checkTerm(tq)
		if err != nil {
			return err
		}

		if strings.ToLower(tq.Current()) == AS {
			err = checkName(tq.Next())
			if err != nil {
				return newError(ErrIllAlias, tq)
			}
		}

		if strings.ToLower(tq.Current()) == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return nil
}
