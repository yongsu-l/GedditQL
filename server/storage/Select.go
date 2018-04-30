package storage

import "GedditQL/server/interpreter"

// Select returns the table with the columns specified
func (db *Database) Select(opts interpreter.SelectOptions) (Response, error) {

	t := Response{}

	if opts.All {
		// ignore columnrefs
		for k, v := range db.Tables[opts.TableRefs[0]].Rows {
			// log.Println(k, v)
			t.Names = append(t.Names, k)
			t.DataTypes = append(t.DataTypes, v.DataType)
			for k, v := range v.Data {
				if len(t.Data) <= k {
					var empty []string
					t.Data = append(t.Data, empty)
					t.Data[k] = append(t.Data[k], v)
				} else {
					t.Data[k] = append(t.Data[k], v)
				}
			}
		}
	} else {
		// Respect columnrefs

		tblName := opts.TableRefs[0]

		for _, v := range opts.ColumnRefs {
			if _, exists := db.Tables[tblName].Rows[v]; exists {
				// Change name if exists in opts.As map
				if _, exists := opts.As[v]; exists {
					t.Names = append(t.Names, opts.As[v])
				} else {
					t.Names = append(t.Names, v)
				}
				t.DataTypes = append(t.DataTypes, db.Tables[tblName].Rows[v].DataType)
				for k, v := range db.Tables[tblName].Rows[v].Data {
					if len(t.Data) <= k {
						var empty []string
						t.Data = append(t.Data, empty)
						t.Data[k] = append(t.Data[k], v)
					} else {
						t.Data[k] = append(t.Data[k], v)
					}
				}
			} else {
				return Response{}, &errorString{"Column ref doesn't exist"}
			}
		}
	}

	if opts.Distinct {
		t.Distinct()
	}

	if opts.Limit > 0 {
		t.Data = t.Data[opts.Limit:]
	}

	// log.Println(t.Data[1][0])

	return t, nil
}
