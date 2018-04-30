package storage

import (
	"GedditQL/server/interpreter"
	"log"
)

// Update updates values where specified in db
func (db *Database) Update(opts *interpreter.UpdateOptions) (Response, error) {
	t := Response{}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// log.Print(opts)

	if _, exists := db.Tables[opts.TableRef]; exists {

		// Update only if all requirements are passed
		tmpTbl := db.Tables[opts.TableRef]

		for i := 0; i < tmpTbl.Length; i++ {
			tmp := make(map[string]string)

			for k, v := range tmpTbl.Rows {
				tmp[k] = v.Data[i]
			}

			if chk, err := opts.Condition(tmp); err != nil {
				return Response{}, &errorString{"Error checking row"}
			} else if chk {
				// Update value at tbl if true
				for k, v := range opts.ValueMap {
					tmpTbl.Rows[k].Data[i] = v
					tmpTbl.Rows[k].DataType = opts.TypeMap[v]
				}
			} else {
				log.Print(chk)
			}
		}

	} else {
		return Response{}, &errorString{"Table doesn't exist in DB"}
	}

	db.Save()

	return t, nil
}
