package interpreter

import (
	"strconv"
)

// CheckSyntax : Syntax Checker
///////////////////////////////////////////////////////////////////////////////////////////////
func CheckSyntax(tks []string) error {
	tq := newQueue(tks)

	switch tk := tq.Current(); tk {
	case EOQ:
		return newError(ErrNoQuery, tq)
	case SELECT:
		return checkSelect(tq.Next())
	case INSERT:
		return checkInsert(tq.Next())
	case UPDATE:
		return checkUpdate(tq.Next())
	case CREATE:
		return checkCreate(tq.Next())
	case DELETE:
		return checkDelete(tq.Next())
	default:
		return newError(ErrInvalQuery, tq)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
func checkSelect(tq *queue) error {
	var err error
	if tq.Current() == EOQ {
		return newError(ErrIllEOQuery, tq)
	}
	if tq.Current() == DISTINCT {
		tq.Next()
	}
	if tq.Current() == ALL {
		tq.Next()
	} else {
		err = checkSelectExprs(tq)
		if err != nil {
			return err
		}
	}

	if tq.Current() == FROM {
		err = checkTableRefs(tq.Next())
		if err != nil {
			return err
		}

		if tq.Current() == WHERE {
			err = checkCondition(tq.Next())
		}

		if err != nil {
			return err
		}

		if tq.Current() == ORDER {
			if tq.Next().Current() != BY {
				return newError(ErrExpToken+"\"BY\"", tq)
			}

			err = checkColumnRef(tq.Next())
			if err != nil {
				return err
			}

			if tq.Current() == ASC || tq.Current() == DESC {
				tq.Next()
			}
		}

		if tq.Current() == LIMIT {
			if _, err := strconv.Atoi(tq.Next().Current()); err != nil {
				return newError(ErrExpInt, tq)
			}
			tq.Next()
		}
	}

	if tq.Current() != EOQ {
		return newError(ErrUnexpToken, tq)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkSelectExprs(tq *queue) error {
	if tq.Current() == EOQ {
		return newError(ErrIllEOQuery, tq)
	}
	if tq.Current() == FROM {
		return newError(ErrNoSelExp, tq)
	}

	for {
		err := checkTerm(tq)
		if err != nil {
			return err
		}

		if tq.Current() == AS {
			err = checkName(tq.Next())
			if err != nil {
				return newError(ErrIllAlias, tq)
			}
		}

		if tq.Current() == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkTerm(tq *queue) error {
	if tq.Current() == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	b, err := checkValue(tq)

	if err != nil {
		return err
	}
	if b {
		return nil
	}

	return checkColumnRef(tq)
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkValue(tq *queue) (bool, error) {
	if tq.Current() == EOQ {
		return false, newError(ErrIllEOQuery, tq)
	}

	if _, err := strconv.Atoi(tq.Current()); err == nil {
		if tq.Next().Current() != "." {
			return true, nil
		}

		if _, err := strconv.Atoi(tq.Next().Current()); err == nil {
			tq.Next()
			return true, nil
		}

		return false, newError(ErrUnexpToken, tq)
	}

	if tq.Current()[0] == '\'' {
		if tq.Current()[len(tq.Current())-1] == '\'' {
			tq.Next()
			return true, nil
		}

		return false, newError(ErrOpenQuote, tq)
	}

	if tq.Current()[0] == '"' {
		if tq.Current()[len(tq.Current())-1] == '"' {
			tq.Next()
			return true, nil
		}
		return false, newError(ErrOpenQuote, tq)
	}

	if tq.Current() == TRUE || tq.Current() == FALSE {
		tq.Next()
		return true, nil
	}

	return false, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkColumnRef(tq *queue) error {
	if tq.Current() == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	err := checkName(tq)
	if err != nil {
		return newError(ErrIllColRef, tq)
	}

	if tq.Current() == DOT {
		if tq.Next().Current() == EOQ {
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

///////////////////////////////////////////////////////////////////////////////////////////////
func checkTableRefs(tq *queue) error {
	if tq.Current() == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	for {
		err := checkName(tq)
		if err != nil {
			return newError(ErrIllTabRef, tq)
		}

		if tq.Current() == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkName(tq *queue) error {
	dt := RESERVED

	for len(dt) > 0 {
		if tq.Current() == dt[0] {
			return newError("", tq)
		}
		dt = dt[1:]
	}

	tq.Next()
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkCondition(tq *queue) error {
	if tq.Current() == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	err := checkExpr(tq)
	if err != nil {
		return err
	}

	if tq.Current() == AND || tq.Current() == OR {
		err := checkExpr(tq.Next())
		if err != nil {
			return err
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkExpr(tq *queue) error {
	if tq.Current() == EOQ {
		return newError(ErrIllEOQuery, tq)
	}
	err := checkTerm(tq)
	if err != nil {
		return err
	}

	switch tk := tq.Current(); tk {
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

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
func checkCreate(tq *queue) error {
	if tq.Current() != TABLE && tq.Current() != VIEW {
		return newError(ErrNoCreate, tq)
	}

	if tq.Next().Current() == IF {
		if tq.Next().Current() != NOT || tq.Next().Current() != EXISTS {
			return newError(ErrUnexpToken, tq)
		}
		tq.Next()
	}

	err := checkName(tq)
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	if tq.Current() != LPAREN {
		return newError(ErrExpToken+"\""+LPAREN+"\"", tq)
	}

	err = checkColumnDefs(tq.Next())
	if err != nil {
		return err
	}

	if tq.Current() != RPAREN {
		return newError(ErrExpToken+"\""+RPAREN+"\"", tq)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkColumnDefs(tq *queue) error {
	if tq.Current() == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	for {
		err := checkName(tq)
		if err != nil {
			return newError(ErrIllColRef, tq)
		}

		err = checkDataType(tq)
		if err != nil {
			return err
		}

	loop:
		for {
			switch tk := tq.Current(); tk {
			case EOQ:
				return newError(ErrIllEOQuery, tq)
			case NOT:
				if tq.Next().Current() != NULL {
					return newError(ErrUnexpToken, tq)
				}
				tq.Next()
				break
			case PRIMARY:
				if tq.Next().Current() != KEY {
					return newError(ErrUnexpToken, tq)
				}
				tq.Next()
				break
			case DEFAULT:
				b, err := checkValue(tq.Next())
				if err != nil {
					return err
				}
				if !b {
					return newError(ErrIllValue, tq)
				}
				break
			// case REFERENCES:
			// 	err = checkName(tq.Next())
			// 	if err != nil {
			// 		return err
			// 	}
			// 	if tq.Current() != "." {
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

		if tq.Current() == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkDataType(tq *queue) error {
	dt := DATATYPES

	for len(dt) > 0 {
		if tq.Current() == dt[0] {
			tq.Next()
			return nil
		}
		dt = dt[1:]
	}

	tq.Next()
	return newError(ErrIllDType, tq)
}

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
func checkInsert(tq *queue) error {
	if tq.Current() != INTO {
		return newError(ErrExpToken+"\"INTO\"", tq)
	}

	err := checkName(tq.Next())
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	if tq.Current() == LPAREN {
		tq.Next()
		for {
			err = checkName(tq)
			if err != nil {
				return newError(ErrIllColRef, tq)
			}

			if tq.Current() == COMMA {
				tq.Next()
			} else {
				break
			}
		}
		if tq.Current() != RPAREN {
			return newError(ErrExpToken+"\""+RPAREN+"\"", tq)
		}
		tq.Next()
	}

	if tq.Current() != VALUES {
		return newError(ErrExpToken+"\"VALUES\"", tq)
	}

	if tq.Next().Current() != LPAREN {
		return newError(ErrExpToken+"\""+LPAREN+"\"", tq)
	}

	tq.Next()
	for {
		b, err := checkValue(tq)
		if err != nil {
			return err
		}
		if !b {
			return newError(ErrIllValue, tq)
		}
		if tq.Current() == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	if tq.Current() != RPAREN {
		return newError(ErrExpToken+"\""+RPAREN+"\"", tq)
	}

	if tq.Next().Current() != EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
func checkUpdate(tq *queue) error {
	if tq.Current() == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	err := checkName(tq)
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	if tq.Current() != SET {
		return newError(ErrExpToken+"\"SET\"", tq)
	}

	err = checkSetExprs(tq.Next())
	if err != nil {
		return err
	}

	if tq.Current() == WHERE {
		err := checkCondition(tq.Next())
		if err != nil {
			return err
		}
	}

	if tq.Current() != EOQ {
		return newError(ErrUnexpToken, tq)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkSetExprs(tq *queue) error {
	if tq.Current() == EOQ {
		return newError(ErrNoSetExp, tq)
	}

	for {
		err := checkColumnRef(tq)
		if err != nil {
			return err
		}

		if tq.Current() != EQ {
			return newError(ErrUnexpToken, tq)
		}

		b, err := checkValue(tq.Next())

		if err != nil {
			return err
		}

		if !b {
			return newError(ErrIllValue, tq)
		}

		if tq.Current() == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
func checkDelete(tq *queue) error {
	if tq.Current() != FROM {
		return newError(ErrUnexpToken, tq)
	}

	err := checkName(tq.Next())
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	if tq.Current() == WHERE {
		err = checkCondition(tq.Next())
		if err != nil {
			return err
		}
	}

	if tq.Current() == LIMIT {
		if _, err := strconv.Atoi(tq.Next().Current()); err != nil {
			return newError(ErrExpInt, tq)
		}
		tq.Next()
	}

	if tq.Current() != EOQ {
		return newError(ErrUnexpToken, tq)
	}

	return nil
}
