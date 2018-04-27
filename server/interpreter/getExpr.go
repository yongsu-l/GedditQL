package interpreter

import "strings"

func getExpr(tq *queue, sa map[string]string) func(map[string]string) (bool, error) {
	var (
		curr1 string
		pos1  int
		curr2 string
		pos2  int
	)

	isVal1, val1, _ := getValue(tq)
	if !isVal1 {
		val1 = getColumnRef(tq)
		curr1 = val1
		pos1 = tq.Pos()

		val, ok := sa[val1]
		if ok {
			val1 = val
		}
	}

	op := strings.ToLower(tq.Current())

	isVal2, val2, _ := getValue(tq.Next())
	if !isVal2 {
		val2 = getColumnRef(tq)
		curr2 = val2
		pos2 = tq.Pos()

		val, ok := sa[val2]
		if ok {
			val2 = val
		}
	}

	return func(dict map[string]string) (bool, error) {
		if !isVal1 {
			val, ok := dict[val1]
			if !ok {
				return false, newError2(ErrNoColRef, curr1, pos1)
			}
			val1 = val
		}

		if !isVal2 {
			val, ok := dict[val2]
			if !ok {
				return false, newError2(ErrNoColRef, curr2, pos2)
			}
			val2 = val
		}

		switch op {
		case "<":
			return val1 < val2, nil
		case "<>", "!=":
			return val1 != val2, nil
		case "<=":
			return val1 <= val2, nil
		case ">":
			return val1 > val2, nil
		case ">=":
			return val1 >= val2, nil
		case "=":
			return val1 == val2, nil
		default:
			return false, nil
		}
	}
}
