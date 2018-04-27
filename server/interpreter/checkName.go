package interpreter

import "strings"

func checkName(tq *queue) error {
	dt := RESERVED

	if dt[strings.ToLower(tq.Current())] {
		return newError("", tq)
	}

	tq.Next()
	return nil
}
