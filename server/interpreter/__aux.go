package interpreter

import (
	"errors"
	"strconv"
)

///////////////////////////////////////////////////////////////////////////////////////////////
type queue struct {
	q []string
	p int
}

func newQueue(lst []string) *queue {
	return &queue{q: lst, p: 1}
}

func (q *queue) Next() *queue {
	q.q = q.q[1:]
	q.p++
	return q
}

func (q *queue) Current() string {
	return q.q[0]
}

func (q *queue) Pos() int {
	return q.p
}

///////////////////////////////////////////////////////////////////////////////////////////////
func newError(msg string, tq *queue) error {
	return errors.New(msg + "  ---> \"" + tq.Current() + "\" at (pos " + strconv.Itoa(tq.Pos()) + ")")
}

func newError2(msg string, current string, pos int) error {
	return errors.New(msg + "  ---> \"" + current + "\" at (pos " + strconv.Itoa(pos) + ")")
}
