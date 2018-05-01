package storage

import "strings"

// Distinct returns the table with only unique columns
func (tbl *Response) Distinct() error {

	e := map[string]bool{}

	// log.Println(tbl.Data)

	for k, v := range tbl.Data {
		// log.Println(k, v)
		// log.Println(strings.Join(v, ""))
		if _, exists := e[strings.Join(v, "")]; exists == false {
			e[strings.Join(v, "")] = true
			// log.Println(tbl.Data[:k])
			// log.Println(tbl.Data[k+1:])
		} else {
			tbl.Data = append(tbl.Data[:k], tbl.Data[k+1:]...)
		}
	}

	// log.Println(tbl.Data)

	return nil
}
