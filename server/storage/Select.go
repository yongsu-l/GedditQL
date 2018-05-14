package storage

import (
	"GedditQL/server/interpreter"
	"log"
	"sort"
)

// Select returns the table with the columns specified
func (db *Database) Select(opts *interpreter.SelectOptions) (Response, error) {

	t := Response{}

	// Check if tableRef is in opts
	if len(opts.TableRefs) > 0 {
		if len(opts.TableRefs) == 1 {
			if _, exists := db.Tables[opts.TableRefs[0]]; exists {
				if opts.All {
					// Ignore columnrefs and loop over the rows
					if opts.Condition != nil {
						for i := 0; i < db.Tables[opts.TableRefs[0]].Length; i++ {
							tmp := make(map[string]string)
							for k, v := range db.Tables[opts.TableRefs[0]].Rows {
								tmp[k] = v.Data[i]
							}
							log.Print(i)

							// Check if the condition function isn't nil
							if chk, err := opts.Condition(tmp); err != nil {
								return Response{}, &errorString{"Error checking row"}
							} else if chk {
								// Append values to the dataset
								for k, v := range db.Tables[opts.TableRefs[0]].Rows {
									// Append the values to the response
									if _, exists := opts.As[k]; exists {
										t.Names = append(t.Names, opts.As[k])
									} else {
										t.Names = append(t.Names, k)
									}
									if len(t.Data) <= i {
										var empty []string
										t.Data = append(t.Data, empty)
										t.Data[i] = append(t.Data[i], v.Data[i])
									} else {
										t.Data[i] = append(t.Data[i], v.Data[i])
									}
								}
								t.RowsAffected++
							}
						}
					} else {
						// If there is no condition, just append all
						for k, v := range db.Tables[opts.TableRefs[0]].Rows {
							if _, exists := opts.As[k]; exists {
								t.Names = append(t.Names, opts.As[k])
							} else {
								t.Names = append(t.Names, k)
							}
							for k, v := range v.Data {
								if len(t.Data) <= k {
									var empty []string
									t.Data = append(t.Data, empty)
									t.Data[k] = append(t.Data[k], v)
								} else {
									t.Data[k] = append(t.Data[k], v)
								}
							}
							t.RowsAffected++
						}
					}
				} else {
					// Respect columnrefs

					tbl := db.Tables[opts.TableRefs[0]]

					if opts.Condition != nil {
						// Check for condition
						for i := 0; i < tbl.Length; i++ {
							tmp := make(map[string]string)
							for k, v := range tbl.Rows {
								tmp[k] = v.Data[i]
							}

							// log.Print(tmp)

							// Check if the condition function isn't nil
							if chk, err := opts.Condition(tmp); err != nil {
								// log.Print("HELLO")
								return Response{}, &errorString{err.Error()}
							} else if chk {
								for _, v := range opts.ColumnRefs {
									// log.Print(v)
									if _, exists := tbl.Rows[v]; exists {
										if _, exists = opts.As[v]; exists {
											t.Names = append(t.Names, opts.As[v])
										} else {
											t.Names = append(t.Names, v)
										}
										if len(t.Data) <= i {
											var empty []string
											t.Data = append(t.Data, empty)
											t.Data[i] = append(t.Data[i], tbl.Rows[v].Data[i])
										} else {
											t.Data[i] = append(t.Data[i], tbl.Rows[v].Data[i])
										}
									} else {
										return Response{}, &errorString{"Column not found"}
									}
								}
								t.RowsAffected++
							}
						}

					} else {

						for _, v := range opts.ColumnRefs {
							if _, exists := db.Tables[opts.TableRefs[0]].Rows[v]; exists {
								// Change name if exists in opts.As map
								if _, exists := opts.As[v]; exists {
									t.Names = append(t.Names, opts.As[v])
								} else {
									t.Names = append(t.Names, v)
								}
								t.DataTypes = append(t.DataTypes, db.Tables[opts.TableRefs[0]].Rows[v].DataType)
								for k, v := range db.Tables[opts.TableRefs[0]].Rows[v].Data {
									if len(t.Data) <= k {
										var empty []string
										t.Data = append(t.Data, empty)
										t.Data[k] = append(t.Data[k], v)
										t.RowsAffected++
									} else {
										t.Data[k] = append(t.Data[k], v)
										t.RowsAffected++
									}
								}
							} else {
								return Response{}, &errorString{"Column ref doesn't exist"}
							}
						}
					}
				}

				if opts.Distinct {
					t.Distinct()
				}

				if opts.Limit > 0 {
					if opts.Limit > len(t.Data) {
						return Response{}, &errorString{"Limit out of range"}
					}
					t.Data = t.Data[:opts.Limit]
				}

				log.Println(len(opts.FuncCols))

				if opts.Order != "" && len(t.Data) > 1 {
					if opts.By == "ASC" || opts.By == "" {
						// Order by asc on default
						for k, v := range t.Names {
							if v == opts.Order {
								sort.SliceStable(t.Data, func(i, j int) bool {
									return t.Data[i][k] < t.Data[j][k]
								})
							}
						}
					} else {
						// Order by desc
						for k, v := range t.Names {
							if v == opts.Order {
								sort.SliceStable(t.Data, func(i, j int) bool {
									return t.Data[i][k] > t.Data[j][k]
								})
							}
						}
					}
				}

				// log.Println(t.Data[1][0])

				return t, nil
			}
		} else {
			// Cartesian product
		}
		// Return error if table doesn't exist
		return Response{}, &errorString{"Table doesn't exist"}
	}
	// Parse the column refs if there are no tables specified
	for _, v := range opts.ColumnRefs {
		var empty []string
		t.Data = append(t.Data, empty)
		t.Data[0] = append(t.Data[0], v)
	}
	return t, nil
}
