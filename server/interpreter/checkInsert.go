package interpreter

import (
	"errors"
	"strings"
)

func checkInsert(tq *queue) error {
	if strings.ToLower(tq.Current()) != INTO {
		return newError(ErrExpToken+"\"INTO\"", tq)
	}

	err := checkName(tq.Next())
	if err != nil {
		return newError(ErrIllTabRef, tq)
	}

	colCount, valCount := 0, 0

	if strings.ToLower(tq.Current()) == LPAREN {
		tq.Next()
		for {
			colCount++
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
		valCount++
		b, _, err := checkValue(tq)
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

	if colCount != 0 && colCount != valCount {
		return errors.New(ErrLenMatch)
	}

	if strings.ToLower(tq.Next().Current()) != EOQ {
		return newError(ErrIllEOQuery, tq)
	}

	return nil
}
