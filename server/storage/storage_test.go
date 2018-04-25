package storage

import (
	"os"
	"testing"
)

var (
	db  *Database
	dir = "test"
)

func TestMain(m *testing.M) {

	// Remove all files from previously failed tests
	os.RemoveAll(dir)

	// Run
	code := m.Run()

	// cleanup
	// os.RemoveAll(dir)

	// Exit
	os.Exit(code)
}

// Test for db shouldn't exist and not exist
func TestNew(t *testing.T) {

	// DB should not exist
	if _, err := os.Stat(dir); err == nil {
		t.Error("Expected no db, but db exists")
	}

	// Create db
	createDB()

	// DB should exist
	if _, err := os.Stat(dir); err != nil {
		t.Error("Expected db, but no db exist")
	}
}

// Test for writing and reading to db
func TestWriteAndRead(t *testing.T) {
	// Write to table first and then write to db after writing to data
	Col1 := Schema{
		Name: "Col1",
		Type: "int",
	}
	Col2 := Schema{
		Name: "Col2",
		Type: "string",
	}

	ColArray := []Schema{Col1, Col2}
	err := db.CreateTable("Test Table", ColArray)

	if err != nil {
		t.Error("Error creating table")
	}

	// log.Println(db.Tables)

	//err := db.Save()
	//if err != nil {
	//	t.Error("Failed to save to file")
	//}

	//err = db.ReadAll()
	//if err != nil {
	//	t.Error("Failed to load file")
	//}

	//log.Println(db.Tables)

}

func TestFrom(t *testing.T) {
	// Make a table within the db then check if it exists
	const tblName = "Test"

	if _, exist := db.Tables[tblName]; exist {
		t.Error("Table shouldn't exist in db")
	}

	// Create table then check if it exists

	var empty []Schema

	db.CreateTable(tblName, empty)

	if _, exist := db.Tables[tblName]; exist == false {
		t.Error("Table should exist in db")
	}

}

func createDB() error {
	var err error
	if db, err = New(dir); err != nil {
		return err
	}

	return nil
}
