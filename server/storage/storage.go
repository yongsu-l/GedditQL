package storage

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
)

const filename = "storage.gob"

// Data stored in row
type Data struct {
}

// Row which will hold all of the data
type Row struct {
	Data []string
}

// Column is the type of columns in the table
type Column struct {
	Name       string
	ColumnType string
}

// Table which holds all of the values in an array as order matters
type Table struct {
	// Name of Table
	Name string
	// Specifies the name of columns
	Columns []Column
	// Data held within each row
	Rows []Row
}

// Database which holds all of the tables
type Database struct {
	Dir    string
	File   string
	Tables map[string]Table
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

// CreateTable with names and type specified
func (db *Database) CreateTable(tablename string, Columns []Column) {
	db.Tables[tablename] = Table{Name: tablename, Columns: Columns}
}

// From returns the table of the tablename
func (db *Database) From(tblName string) Table {
	return db.Tables[tblName]
}

//func (db *Database) Select(Cols []string) Table {
//
//}

/*

Database: Holds tables in map that ties name to struct
Tables: Holds array of row name

*/
