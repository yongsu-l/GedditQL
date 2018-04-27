package interpreter

func getSetExprs(tq *queue) (map[string]string, map[string]string) {
	valueMap := map[string]string{}
	typeMap := map[string]string{}

	for {
		colRef := tq.Current()
		_, val, typ := getValue(tq.Next().Next())
		valueMap[colRef] = val
		typeMap[val] = typ

		if tq.Current() == COMMA {
			tq.Next()
		} else {
			break
		}
	}

	return valueMap, typeMap
}
