package interpreter

import "strings"

func checkDataType(tq *queue) error {
	dt := DATATYPES

	if dt[strings.ToLower(tq.Current())] {
		tq.Next()
		return nil
	}

	return newError(ErrIllDType, tq)
}
