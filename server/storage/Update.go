package storage

import (
	"GedditQL/server/interpreter"
	"errors"
	"fmt"
)

// Update updates values where specified in db
func (db *Database) Update(opts *interpreter.UpdateOptions) (Response, error) {
	t := Response{}

	if _, exists := db.Tables[opts.TableRef]; exists {

		// Update only if all requirements are passed
		tmpTbl := db.Tables[opts.TableRef]

		for i := 0; i < tmpTbl.Length; i++ {
			tmp := make(map[string]string)

			for k, v := range tmpTbl.Rows {
				tmp[k] = v.Data[i]
			}

			// Check if condition is nil
			if opts.Condition != nil {
				if chk, err := opts.Condition(tmp); err != nil {
					db.Load()
					return Response{}, errors.New("Error checking row")
				} else if chk {
					// Update value at tbl if true
					for k, v := range opts.ValueMap {
						// Check if the value exists in db
						if _, exists := tmpTbl.Rows[k]; exists {
							tmpTbl.Rows[k].Data[i] = v
							tmpTbl.Rows[k].DataType = opts.TypeMap[v]
						} else {
							// Reload the db from saved file
							db.Load()
							return Response{}, errors.New("Column does not exist in DB")
						}
					}
				}
			}
		}

	} else {
		db.Load()
		return Response{}, errors.New("Table doesn't exist in DB")
	}

	db.Save()

	t.Result = fmt.Sprintf("Updated table of: %s", opts.TableRef)

	return t, nil
}
