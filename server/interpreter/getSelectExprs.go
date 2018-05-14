package interpreter

import "strings"

func getSelectExprs(tq *queue) ([]string, map[string]string, map[string]string, []string, map[string]string) {
	columnRefs := []string{}
	// columnRef: AS
	as := map[string]string{}
	// AS: columnRef
	sa := map[string]string{}
	// function columns
	fc := []string{}
	// function map
	fm := map[string]string{}

	for {
		if tq.Ind(1) == LPAREN {
			if tq.Ind(2) == ALL {
				f := strings.ToLower(tq.Current())
				s := tq.Next().Next().Current()

				fc = append(fc, s)
				fm[s] = f
			} else {
				f := strings.ToLower(tq.Current())
				s := getColumnRef(tq.Next().Next())

				fc = append(fc, s)
				fm[s] = f

			}
			if tq.Next().Current() == COMMA {
				tq.Next()
				continue
			} else {
				break
			}
		}

		s := getColumnRef(tq)
		columnRefs = append(columnRefs, s)
		// checks for AS clause
		if strings.ToLower(tq.Current()) == AS {
			a := tq.Next().Current()
			as[s] = a
			sa[a] = s
			tq.Next()
		}

		// checks if there are more columnRefs
		if strings.ToLower(tq.Current()) == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return columnRefs, as, sa, fc, fm
}
