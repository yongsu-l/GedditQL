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
		tk := strings.ToLower(tq.Current())
		if tk == SUM || tk == COUNT {
			if strings.ToLower(tq.Ind(1)) == LPAREN {
				err := checkColumnRef(tq.Next().Next())

				if err != nil {
					return newError(ErrIllColRef, tq)
				}

				if strings.ToLower(tq.Current()) == RPAREN {
					if tq.Next().Current() == COMMA {
						tq.Next()
						continue
					} else {
						break
					}
				} else {
					return newError(ErrExpToken+")", tq)
				}
			}

		}

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
