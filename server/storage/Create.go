package storage

import (
	"GedditQL/server/interpreter"
	"log"
)

// Create creates a table specified
func (db *Database) Create(opts *interpreter.CreateOptions) (Response, error) {

	t := Response{}

	if opts.Type == "table" {
		// Implement Table creation
		if opts.IfNotExists {
			// Check if the table exists in the db
			if _, exists := db.Tables[opts.TableRef]; exists {
				return Response{}, &errorString{"Table exists in DB"}
			} else {
				// If the table doesn't exist, create table
				tmpTable := Table{}

				for _, v := range opts.ColumnDefs {
					// Insert into tmpTable
					tmpTable.Rows[v.Name] = &Data{DataType: v.DataType}
				}

				db.Tables[opts.TableRef] = &tmpTable
				db.Save()
			}
		} else {
			emptyRow := make(map[string]*Data)
			db.Tables[opts.TableRef] = &Table{Rows: emptyRow, Length: 0}

			for _, v := range opts.ColumnDefs {
				db.Tables[opts.TableRef].Rows[v.Name] = &Data{DataType: v.DataType}
			}

			db.Save()
			log.Print(db.Tables[opts.TableRef])
			t.RowsAffected = 1
		}
	} else if opts.Type == "view" {
		// Implement Database creation
	}

	return t, nil
}
