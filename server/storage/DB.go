package storage

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const filename = "storage.gob"

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
		Tables: make(map[string]*Table),
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
	db.Tables[tblName] = &Table{Rows: Rows}
	// Save on each CreateTable
	return db.Save()
}

// ColumnNames ///
func (db *Database) ColumnNames(tblName string) []string {

	var columns []string

	log.Println(db.Tables[tblName])

	if db.Tables[tblName] != nil {
		for k := range db.Tables[tblName].Rows {
			columns = append(columns, k)
		}
	}

	return columns
}

// Count returns the count of rows in table
func (db *Database) Count(tblName string) int {
	return len(db.Tables[tblName].Rows)
}

// Distinct returns the table with only unique columns
func (tbl *Response) Distinct() error {

	e := map[string]bool{}

	// log.Println(tbl.Data)

	for k, v := range tbl.Data {
		// log.Println(k, v)
		// log.Println(strings.Join(v, ""))
		if _, exists := e[strings.Join(v, "")]; exists == false {
			e[strings.Join(v, "")] = true
			// log.Println(tbl.Data[:k])
			// log.Println(tbl.Data[k+1:])
		} else {
			tbl.Data = append(tbl.Data[:k], tbl.Data[k+1:]...)
		}
	}

	// log.Println(tbl.Data)

	return nil
}

// Update updates the value at the table specified

// CreateTable with names and type specified
//func (db *Database) CreateTable(tablename string, Columns []Column) {
//	db.Tables[tablename] = Table{Columns: Columns}
//}

/*

Database: Holds tables in map that ties name to struct
Tables: Holds array of row name

*/