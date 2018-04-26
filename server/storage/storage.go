package storage

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
)

const filename = "storage.gob"

// Data which will hold all of the data
type Data struct {
	Type string
	Data []string
}

// Table which holds all of the values in an array as order matters
type Table struct {
	// Data held within each row linearly
	Rows map[string]*Data
}

// ReturnTable is the type which is first an array of column names, then an array of types, then an array of all the data
type ReturnTable struct {
	Names []string
	Types []string
	Data  [][]string
}

// Database which holds all of the tables
type Database struct {
	Dir    string
	File   string
	Tables map[string]Table
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// New Initializes a struct of DB and creates file and directory if needed
func New(dir string) (*Database, error) {

	// Clean the directory that we will be working in
	dir = filepath.Clean(dir)

	// The fullpath of the directory that we will place the file in
	db := Database{
		Dir:    dir,
		File:   filename,
		Tables: make(map[string]Table),
	}

	// Create db if exists and create db if it doesn't exist
	if _, err := os.Stat(dir); err == nil {
		log.Printf("DB already exists in dir: %s", dir)
		// Then check if the filename is already in the directory
		if _, err := os.Stat(filepath.Join(dir, filename)); err == nil {
			log.Printf("File already exists in dir: %s", filepath.Join(dir, filename))
			return &db, err
		}
		// Create file if file doesn't exist
		log.Printf("Creating file at %s", filepath.Join(dir, filename))
		file, err := os.Create(filepath.Join(dir, filename))
		if err != nil {
			log.Fatal("Error creating file: ", err.Error())
		}
		file.Close()
		return &db, err

	}
	// Create directory if directory doesn't exist
	log.Printf("Creating filepath at %s", dir)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal("Error creating path")
	}
	// If dir doesn't exist, create the file as well
	log.Printf("Creating file at %s", filepath.Join(dir, filename))
	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		log.Fatal("Error creating file: ", err.Error())
	}
	file.Close()
	return &db, err

}

// ReadAll from file
func (db *Database) ReadAll() error {
	file, err := os.Open(filepath.Join(db.Dir, db.File))
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(db)
	}
	file.Close()
	return err
}

// Save to file
func (db *Database) Save() error {

	// Assume that you can't run Write without New and file is already created
	file, err := os.Create(filepath.Join(db.Dir, db.File))
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(db)
	}
	file.Close()
	return err

}

// CheckData checks if the column name specified exists in table
func (db *Database) CheckData(tblName string, val *Data) bool {

	return true
}

// CreateTable creates a table in the current database
func (db *Database) CreateTable(tblName string, Rows map[string]*Data) error {
	db.Tables[tblName] = Table{Rows: Rows}
	// Save on each CreateTable
	return db.Save()
}

// ColumnNames ///
func (db *Database) ColumnNames(tblName string) []string {
	var columns []string

	for k := range db.Tables[tblName].Rows {
		columns = append(columns, k)
	}

	return columns
}

// Count returns the count of rows in table
func (db *Database) Count(tblName string) int {
	return len(db.Tables[tblName].Rows)
}

// Normalize makes it so that if there are no values that were inserted into the db, it will insert an nil value appropriate
func (db *Database) Normalize(tblName string, length int) error {

	for k := range db.Tables[tblName].Rows {
		if len(db.Tables[tblName].Rows[k].Data) < length {
			db.Tables[tblName].Rows[k].Data = append(db.Tables[tblName].Rows[k].Data, "")
		}
	}

	return nil
}

// InsertInto inserts into the table with the tblName specified along with the column and value pair in the map.
func (db *Database) InsertInto(tblName string, insertion map[string]string) error {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var length int

	// Check if tbl is in db
	if _, ok := db.Tables[tblName]; ok {
		// tbl := db.Tables[tblName]
		for k, v := range insertion {
			// Check if column name exists in table, otherwise throw error
			if _, ok = db.Tables[tblName].Rows[k]; ok {

				// Append new value at the end
				db.Tables[tblName].Rows[k].Data = append(db.Tables[tblName].Rows[k].Data, v)
				// log.Println(db.Tables[tblName].Rows[k].Data)
				if len(db.Tables[tblName].Rows[k].Data) > length {
					length = len(db.Tables[tblName].Rows[k].Data)
				}

			} else {
				return &errorString{"Column name not in database"}
			}
		}
	} else {
		return &errorString{"Table not in database"}
	}

	// Normalize all values so that length are all equal among the arrays
	db.Normalize(tblName, length)

	// for k, v := range db.Tables[tblName].Rows {
	// 	log.Println(k, v.Data)
	// }

	return db.Save()
}

// Select returns the table with the columns specified
func (db *Database) Select(colNames []string, tblName string) (ReturnTable, error) {
	t := ReturnTable{}

	for _, v := range colNames {
		if _, exist := db.Tables[tblName].Rows[v]; exist {
			t.Names = append(t.Names, v)
			t.Types = append(t.Types, db.Tables[tblName].Rows[v].Type)
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
			return ReturnTable{}, &errorString{"Column name doesn't exist"}
		}
	}

	return t, nil
}

// CreateTable with names and type specified
//func (db *Database) CreateTable(tablename string, Columns []Column) {
//	db.Tables[tablename] = Table{Columns: Columns}
//}

/*

Database: Holds tables in map that ties name to struct
Tables: Holds array of row name

*/
