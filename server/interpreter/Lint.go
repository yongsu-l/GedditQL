package interpreter

import (
	"strings"
)

// Lint : Syntax Checker
///////////////////////////////////////////////////////////////////////////////////////////////
func Lint(tks []string) error {
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
