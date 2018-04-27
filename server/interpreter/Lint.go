package interpreter

import (
	"strconv"
	"strings"
)

// CheckSyntax : Syntax Checker
///////////////////////////////////////////////////////////////////////////////////////////////
func CheckSyntax(tks []string) error {
	tq := newQueue(tks)

	switch tk := strings.ToLower(tq.Current()); tk {
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

///////////////////////////////////////////////////////////////////////////////////////////////
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

///////////////////////////////////////////////////////////////////////////////////////////////
func checkTerm(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
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

func checkNumeric(tq *queue) (bool, error) {
	if _, err := strconv.Atoi(strings.ToLower(tq.Current())); err == nil {
		if strings.ToLower(tq.Next().Current()) != "." {
			return true, nil
		}

		if _, err := strconv.Atoi(strings.ToLower(tq.Next().Current())); err == nil {
			tq.Next()
			return true, nil
		}

		return false, newError(ErrIllValue, tq)
	}
	return false, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkValue(tq *queue) (bool, error) {
	if strings.ToLower(tq.Current()) == EOQ {
		return false, newError(ErrIllEOQuery, tq)
	}

	if strings.ToLower(tq.Current()) == SUB {
		b, err := checkNumeric(tq.Next())
		if !b || err != nil {
			return false, newError(ErrIllValue, tq)
		}
		return true, nil
	}

	b, err := checkNumeric(tq)
	if err != nil {
		return false, err
	}
	if b {
		return true, nil
	}

	if strings.ToLower(tq.Current())[0] == '\'' {
		if strings.ToLower(tq.Current())[len(strings.ToLower(tq.Current()))-1] == '\'' {
			tq.Next()
			return true, nil
		}

		return false, newError(ErrOpenQuote, tq)
	}

	if strings.ToLower(tq.Current())[0] == '"' {
		if strings.ToLower(tq.Current())[len(strings.ToLower(tq.Current()))-1] == '"' {
			tq.Next()
			return true, nil
		}
		return false, newError(ErrOpenQuote, tq)
	}

	if strings.ToLower(tq.Current()) == TRUE || strings.ToLower(tq.Current()) == FALSE {
		tq.Next()
		return true, nil
	}

	return false, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
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

///////////////////////////////////////////////////////////////////////////////////////////////
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

///////////////////////////////////////////////////////////////////////////////////////////////
func checkName(tq *queue) error {
	dt := RESERVED

	if dt[strings.ToLower(tq.Current())] {
		return newError("", tq)
	}

	tq.Next()
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
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

///////////////////////////////////////////////////////////////////////////////////////////////
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

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
func checkCreate(tq *queue) error {
	if strings.ToLower(tq.Current()) != TABLE && strings.ToLower(tq.Current()) != VIEW {
		return newError(ErrNoCreate, tq)
	}

	if strings.ToLower(tq.Next().Current()) == IF {
		if strings.ToLower(tq.Next().Current()) != NOT || strings.ToLower(tq.Next().Current()) != EXISTS {
			return newError(ErrUnexpToken, tq)
		}
		tq.Next()
	}

	err := checkName(tq)
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	if strings.ToLower(tq.Current()) != LPAREN {
		return newError(ErrExpToken+"\""+LPAREN+"\"", tq)
	}

	err = checkColumnDefs(tq.Next())
	if err != nil {
		return err
	}

	if strings.ToLower(tq.Current()) != RPAREN {
		return newError(ErrExpToken+"\""+RPAREN+"\"", tq)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
func checkColumnDefs(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
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

///////////////////////////////////////////////////////////////////////////////////////////////
func checkDataType(tq *queue) error {
	dt := DATATYPES

	if dt[strings.ToLower(tq.Current())] {
		tq.Next()
		return nil
	}

	tq.Next()
	return newError(ErrIllDType, tq)
}

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
func checkInsert(tq *queue) error {
	if strings.ToLower(tq.Current()) != INTO {
		return newError(ErrExpToken+"\"INTO\"", tq)
	}

	err := checkName(tq.Next())
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	if strings.ToLower(tq.Current()) == LPAREN {
		tq.Next()
		for {
			err = checkName(tq)
			if err != nil {
				return newError(ErrIllColRef, tq)
			}

			if strings.ToLower(tq.Current()) == COMMA {
				tq.Next()
			} else {
				break
			}
		}
		if strings.ToLower(tq.Current()) != RPAREN {
			return newError(ErrExpToken+"\""+RPAREN+"\"", tq)
		}
		tq.Next()
	}

	if strings.ToLower(tq.Current()) != VALUES {
		return newError(ErrExpToken+"\"VALUES\"", tq)
	}

	if strings.ToLower(tq.Next().Current()) != LPAREN {
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
		if strings.ToLower(tq.Current()) == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	if strings.ToLower(tq.Current()) != RPAREN {
		return newError(ErrExpToken+"\""+RPAREN+"\"", tq)
	}

	if strings.ToLower(tq.Next().Current()) != EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
func checkUpdate(tq *queue) error {
	if strings.ToLower(tq.Current()) == EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	err := checkName(tq)
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	if strings.ToLower(tq.Current()) != SET {
		return newError(ErrExpToken+"\"SET\"", tq)
	}

	err = checkSetExprs(tq.Next())
	if err != nil {
		return err
	}

	if strings.ToLower(tq.Current()) == WHERE {
		err := checkCondition(tq.Next())
		if err != nil {
			return err
		}
	}

	if strings.ToLower(tq.Current()) != EOQ {
		return newError(ErrUnexpToken, tq)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
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

		b, err := checkValue(tq.Next())

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

///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////
func checkDelete(tq *queue) error {
	if strings.ToLower(tq.Current()) != FROM {
		return newError(ErrUnexpToken, tq)
	}

	err := checkName(tq.Next())
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	if strings.ToLower(tq.Current()) == WHERE {
		err = checkCondition(tq.Next())
		if err != nil {
			return err
		}
	}

	if strings.ToLower(tq.Current()) == LIMIT {
		if _, err := strconv.Atoi(strings.ToLower(tq.Next().Current())); err != nil {
			return newError(ErrExpInt, tq)
		}
		tq.Next()
	}

	if strings.ToLower(tq.Current()) != EOQ {
		return newError(ErrUnexpToken, tq)
	}

	return nil
}
