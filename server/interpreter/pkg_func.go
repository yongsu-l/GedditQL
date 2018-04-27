package interpreter

import (
	"errors"
	"strconv"
)

///////////////////////////////////////////////////////////////////////////////////////////////
func newError(msg string, tq *queue) error {
	return errors.New(msg + "  ---> \"" + tq.Current() + "\" at (pos " + strconv.Itoa(tq.Pos()) + ")")
}

///////////////////////////////////////////////////////////////////////////////////////////////
func newError2(msg string, current string, pos int) error {
	return errors.New(msg + "  ---> \"" + current + "\" at (pos " + strconv.Itoa(pos) + ")")
}
