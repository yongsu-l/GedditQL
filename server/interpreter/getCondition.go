package interpreter

import "strings"

func getCondition(tq *queue, sa map[string]string) func(map[string]string) (bool, error) {
	res := getExpr(tq, sa)

	for strings.ToLower(tq.Current()) == AND || strings.ToLower(tq.Current()) == OR {
		op := strings.ToLower(tq.Current())
		next := getExpr(tq.Next(), sa)
		prev := res

		res = func(dict map[string]string) (bool, error) {
			b1, err := prev(dict)
			if err != nil {
				return false, err
			}
			b2, err := next(dict)
			if err != nil {
				return false, err
			}

			switch op {
			case AND:
				return b1 && b2, nil
			case OR:
				return b1 || b2, nil
			default:
				return false, nil
			}
		}
	}

	return res
}
