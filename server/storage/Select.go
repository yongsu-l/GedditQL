package storage

import (
	"GedditQL/server/interpreter"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
)

// Select returns the table with the columns specified
func (db *Database) Select(opts *interpreter.SelectOptions) (Response, error) {

	t := Response{}

	// Check if tableRef is in opts
	if len(opts.TableRefs) > 0 {
		if len(opts.TableRefs) == 1 {
			if _, exists := db.Tables[opts.TableRefs[0]]; exists {

				tbl := db.Tables[opts.TableRefs[0]]

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
						}
					}
				} else {
					// Respect columnrefs

					if opts.Condition != nil {
						// Check for condition
						for i := 0; i < tbl.Length; i++ {
							tmp := make(map[string]string)
							for k, v := range tbl.Rows {
								tmp[k] = v.Data[i]
							}

							// Check for condition
							if chk, err := opts.Condition(tmp); err != nil {
								return Response{}, err
							} else if chk {
								for _, v := range opts.ColumnRefs {

									// Check if column exists
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
							}
						}

					} else {

						for _, v := range opts.ColumnRefs {
							if _, exists := tbl.Rows[v]; exists {
								// Change name if exists in opts.As map
								if _, exists := opts.As[v]; exists {
									t.Names = append(t.Names, opts.As[v])
								} else {
									t.Names = append(t.Names, v)
								}
								t.DataTypes = append(t.DataTypes, tbl.Rows[v].DataType)
								for k, v := range tbl.Rows[v].Data {
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
				}

				if opts.Distinct {
					t.Distinct()
				}

				if len(opts.FuncCols) != 0 {
					err := db.funcInterpret(&t, opts)
					if err != nil {
						return Response{}, err
					}
				}

				if opts.Limit > 0 {
					if opts.Limit > len(t.Data) {
						return Response{}, errors.New("Limit out of range")
					}
					t.Data = t.Data[:opts.Limit]
				}

				t.RowsAffected = len(t.Data)

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
			} else {
				// Return error if table doesn't exist
				return Response{}, errors.New("Table does not exist")
			}
		} else {
			// Cartesian product
			return Response{}, errors.New("Cartersian Product not implemented")
		}

		return t, nil
	}
	// Parse the column refs if there are no tables specified
	for _, v := range opts.ColumnRefs {
		var empty []string
		t.Data = append(t.Data, empty)
		t.Data[0] = append(t.Data[0], v)
	}
	return t, nil
}

func (db *Database) funcInterpret(res *Response, opts *interpreter.SelectOptions) error {

	var tmp []string
	var tmpColName []string

	for _, v := range opts.FuncCols {
		// Check if sum or count
		if opts.FuncMap[v] == "sum" {

			tmpColName = append(tmpColName, "SUM")

			// Check if the column exists in table
			tbl := db.Tables[opts.TableRefs[0]]
			if _, exists := tbl.Rows[v]; exists {
				// Check if we can sum the column
				if tbl.Rows[v].DataType == "string" {
					return errors.New("Cannot sum string")
				} else if tbl.Rows[v].DataType == "boolean" {
					return errors.New("Cannot sum boolean")
				} else {
					// Sum either float or int
					// Check if float or int
					if tbl.Rows[v].DataType == "float" {
						tmpSum := 0.0
						for _, v := range tbl.Rows[v].Data {
							if f, err := strconv.ParseFloat(v, 64); err != nil {
								return err
							} else {
								tmpSum += f
							}
						}
						tmp = append(tmp, fmt.Sprint(tmpSum))
					} else if tbl.Rows[v].DataType == "int" {
						tmpSum := 0
						for _, v := range tbl.Rows[v].Data {
							if i, err := strconv.Atoi(v); err != nil {
								return err
							} else {
								tmpSum += i
							}
						}
						tmp = append(tmp, fmt.Sprint(tmpSum))
					} else {
						// Not yet implemented
					}
				}
			} else {
				return errors.New("Column doesn't exist in DB")
			}
		} else if opts.FuncMap[v] == "count" {

			tmpColName = append(tmpColName, "COUNT")
			tmpCount := 0
			// Append the current length of the data
			if opts.Condition != nil {

				for i := 0; i < db.Tables[opts.TableRefs[0]].Length; i++ {
					tmp := make(map[string]string)
					for k, v := range db.Tables[opts.TableRefs[0]].Rows {
						tmp[k] = v.Data[i]
					}
					log.Print(i)

					// Check if the condition function isn't nil
					if chk, err := opts.Condition(tmp); err != nil {
						return errors.New("Error checking column")
					} else if chk {
						tmpCount++
					}
				}
			} else {
				tmpCount = db.Tables[opts.TableRefs[0]].Length
			}
			tmp = append(tmp, fmt.Sprint(tmpCount))

		} else {
			// There is nothing to do so do nothing
			return errors.New("Function not yet implemented")
		}
	}

	if len(tmpColName) != 0 {
		res.Names = append(res.Names, tmpColName...)
		if len(res.Data) == 0 {
			var empty []string
			res.Data = append(res.Data, empty)
			res.Data[0] = append(res.Data[0], tmp...)
		} else {
			for k := range res.Data {
				res.Data[k] = append(res.Data[k], tmp...)
			}
		}
	}

	return nil
}
