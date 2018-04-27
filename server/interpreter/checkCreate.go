package interpreter

import "strings"

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
		return newError(ErrNoColDef, tq)
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
