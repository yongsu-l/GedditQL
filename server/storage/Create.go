package storage

import "GedditQL/server/interpreter"

// Create creates a table specified
func (db *Database) Create(opts interpreter.CreateOptions) (Response, error) {

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
			tmpTable := Table{}

			for _, v := range opts.ColumnDefs {
				tmpTable.Rows[v.Name] = &Data{DataType: v.DataType}
			}

			db.Tables[opts.TableRef] = &tmpTable
			db.Save()
		}
	} else if opts.Type == "database" {
		// Implement Database creation
	}

	return t, nil
}
