package storage

import "GedditQL/server/interpreter"

// Delete deletes from the table where the condition is specified
func (db *Database) Delete(opts *interpreter.DeleteOptions) (Response, error) {
	t := Response{}

	if _, exists := db.Tables[opts.TableRef]; exists {
		// If table exist, loop through db until limit if specified
		if opts.Limit > 0 {
			// Delete only on certain areas
			for i := 0; i < opts.Limit; i++ {
				if i > db.Tables[opts.TableRef].Length {
					// If the limit exceeded, break away from for loop
					break
				} else {
					// column name mapped to value of each column
					tmp := make(map[string]string)

					for k, v := range db.Tables[opts.TableRef].Rows {
						tmp[k] = v.Data[i]
					}

					// Check against condition
					if chk, err := opts.Condition(tmp); err != nil {
						return Response{}, &errorString{"Error checking row"}
					} else if chk {
						// Delete value at row if true
						for k := range tmp {
							db.Tables[opts.TableRef].Rows[k].Data = append(db.Tables[opts.TableRef].Rows[k].Data[:i], db.Tables[opts.TableRef].Rows[k].Data[i+1:]...)
						}
						db.Tables[opts.TableRef].Length = db.Tables[opts.TableRef].Length - 1
					}
				}
			}
		} else {
			// Delete on all occurence
			for i := 0; i < db.Tables[opts.TableRef].Length; i++ {
				tmp := make(map[string]string)

				for k, v := range db.Tables[opts.TableRef].Rows {
					tmp[k] = v.Data[i]
				}

				// Check against condition
				if chk, err := opts.Condition(tmp); err != nil {
					return Response{}, &errorString{"Error checking row"}
				} else if chk {
					// Delete value at row if true
					for k := range tmp {
						db.Tables[opts.TableRef].Rows[k].Data = append(db.Tables[opts.TableRef].Rows[k].Data[:i], db.Tables[opts.TableRef].Rows[k].Data[i+1:]...)
					}
					db.Tables[opts.TableRef].Length = db.Tables[opts.TableRef].Length - 1
				}
			}
		}
	} else {
		return Response{}, &errorString{"Table doesn't exist in db"}
	}

	db.Save()

	return t, nil
}
