package interpreter

import "strings"

func getSelectExprs(tq *queue) ([]string, map[string]string, map[string]string) {
	columnRefs := []string{}
	// columnRef: AS
	as := map[string]string{}
	// AS: columnRef
	sa := map[string]string{}

	for {
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

	return columnRefs, as, sa
}
