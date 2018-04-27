package interpreter

import "strings"

func checkColumnDefs(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	for {
		err := checkName(tq)
		if err != nil {
			return newError(ErrIllColRef, tq)
		}

		dType := tq.Current()
		err = checkDataType(tq)
		if err != nil {
			return err
		}

	loop:
		for {
			switch tk := strings.ToLower(tq.Current()); tk {
			case EOQ:
				return newError(ErrIllEOQuery, tq)
			case NOT:
				if strings.ToLower(tq.Next().Current()) != NULL {
					return newError(ErrUnexpToken, tq)
				}
				tq.Next()
				break
			case PRIMARY:
				if strings.ToLower(tq.Next().Current()) != KEY {
					return newError(ErrUnexpToken, tq)
				}
				tq.Next()
				break
			case DEFAULT:
				val := tq.Next().Current()
				pos := tq.Pos()
				b, typ, err := checkValue(tq)
				if err != nil {
					return err
				}
				if !b {
					return newError(ErrIllValue, tq)
				}
				if dType != typ {
					return newError2(ErrValMatch, val, pos)
				}
				break
			// case REFERENCES:
			// 	err = checkName(tq.Next())
			// 	if err != nil {
			// 		return err
			// 	}
			// 	if strings.ToLower(tq.Current()) != "." {
			// 		return newError(ErrUnexpToken, tq)
			// 	}
			// 	err = checkName(tq.Next())
			// 	if err != nil {
			// 		return err
			// 	}
			case NULL, UNIQUE:
				tq.Next()
				break
			default:
				break loop
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
