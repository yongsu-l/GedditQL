package storage

import (
	"GedditQL/server/interpreter"
	// "log"
	"fmt"
)

// Insert inserts into the table specified
func (db *Database) Insert(opts *interpreter.InsertOptions) (Response, error) {

	t := Response{}
	var length int

	if _, exists := db.Tables[opts.TableRef]; exists {
		// Only insert if the Table exists
		for k, v := range opts.ColumnRefs {
			// Check if column name exists in table
			if _, exists := db.Tables[opts.TableRef].Rows[v]; exists {

				// Append new value at the end
				db.Tables[opts.TableRef].Rows[v].Data = append(db.Tables[opts.TableRef].Rows[v].Data, opts.Values[k])
				// Measure length for normalize at end of loop
				if len(db.Tables[opts.TableRef].Rows[v].Data) > length {
					length = len(db.Tables[opts.TableRef].Rows[v].Data)
				}
			} else {
				return Response{}, &errorString{"Column name not in database"}
			}
		}
	} else {
		// Spit error if table doesn't exist in db
		return Response{}, &errorString{"Table doesn't exist in db"}
	}

	db.normalize(opts.TableRef, length)

	// log.Print(db.Tables[opts.TableRef].Rows["col1"].Data)
	db.Tables[opts.TableRef].Length = length

	db.Save()

	t.Result = fmt.Sprintf("Inserted values into: %s", opts.TableRef)

	return t, nil
}

// Normalize makes it so that if there are no values that were inserted into the db, it will insert an nil value appropriate
func (db *Database) normalize(tblName string, length int) error {

	for k := range db.Tables[tblName].Rows {
		if len(db.Tables[tblName].Rows[k].Data) < length {
			db.Tables[tblName].Rows[k].Data = append(db.Tables[tblName].Rows[k].Data, "")
		}
	}

	return nil
}
