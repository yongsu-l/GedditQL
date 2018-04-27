package interpreter

import (
	"strconv"
	"strings"
)

func checkSelect(tq *queue) error {
	var err error
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrIllEOQuery, tq)
	}
	if strings.ToLower(tq.Current()) == DISTINCT {
		tq.Next()
	}
	if strings.ToLower(tq.Current()) == ALL {
		tq.Next()
	} else {
		err = checkSelectExprs(tq)
		if err != nil {
			return err
		}
	}

	if strings.ToLower(tq.Current()) == FROM {
		err = checkTableRefs(tq.Next())
		if err != nil {
			return err
		}

		if strings.ToLower(tq.Current()) == WHERE {
			err = checkCondition(tq.Next())
		}

		if err != nil {
			return err
		}

		if strings.ToLower(tq.Current()) == ORDER {
			if strings.ToLower(tq.Next().Current()) != BY {
				return newError(ErrExpToken+"\"BY\"", tq)
			}

			err = checkColumnRef(tq.Next())
			if err != nil {
				return err
			}

			if strings.ToLower(tq.Current()) == ASC || strings.ToLower(tq.Current()) == DESC {
				tq.Next()
			}
		}

		if strings.ToLower(tq.Current()) == LIMIT {
			if _, err := strconv.Atoi(strings.ToLower(tq.Next().Current())); err != nil {
				return newError(ErrExpInt, tq)
			}
			tq.Next()
		}
	}

	if strings.ToLower(tq.Current()) != EOQ {
		return newError(ErrUnexpToken, tq)
	}

	return nil
}
