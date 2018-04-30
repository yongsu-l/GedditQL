package storage

import (
	"GedditQL/server/interpreter"
	"GedditQL/server/parser"
	"os"
	"testing"
)

var (
	db      *Database
	dir     = "test"
	tblName = "testTbl"
	col1    = "col1"
	col2    = "col2"
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
	data := &Data{DataType: "string"}
	r1 := make(map[string]*Data)
	r1[col1] = data
	r1[col2] = data

	//ColArray := []Schema{Col1, Col2}
	err := db.CreateTable(tblName, r1)

	if err != nil {
		t.Error("Error creating table")
	}

	if _, ok := db.Tables[tblName]; ok == false {
		t.Fatal("New table was not inserted")
	} else if _, ok := db.Tables[tblName].Rows[col1]; ok == false {
		t.Fatal("New coliumn was not inserted")
	}

	// log.Println(db.Tables[tblName].Rows[col1])

	err = db.Save()
	if err != nil {
		t.Error("Failed to save to file")
	}

	err = db.ReadAll()
	if err != nil {
		t.Error("Failed to load file")
	}

	if _, ok := db.Tables[tblName]; ok == false {
		t.Fatal("New table was not inserted in save")
	} else if _, ok := db.Tables[tblName].Rows[col1]; ok == false {
		t.Fatal("New coliumn was not inserted in save")
	}

	// log.Println(db.Tables["Schema1"])

}

func TestFrom(t *testing.T) {
	// Make a table within the db then check if it exists
	const tblName = "Test"

	if _, exist := db.Tables[tblName]; exist {
		t.Error("Table shouldn't exist in db")
	}

	// Create table then check if it exists

	//var empty []Schema

	//db.CreateTable(tblName, empty)

	//if _, exist := db.Tables[tblName]; exist == false {
	//	t.Error("Table should exist in db")
	//}

}

// func TestInsertInto(t *testing.T) {
// 	insertion := make(map[string]string)
//
// 	insertion[col1] = "hello"
// 	insertion[col2] = "world"
//
// 	db.InsertInto(tblName, insertion)
//
// 	insertion = make(map[string]string)
// 	insertion[col1] = "second round"
//
// 	db.InsertInto(tblName, insertion)
// }

func TestInsert(t *testing.T) {
	opts := interpreter.InsertOptions{
		TableRef:   tblName,
		ColumnRefs: []string{col1, col2},
		Values:     []string{"\"hello\"", "\"world\""},
		ValueTypes: []string{"string", "string"},
	}

	db.Insert(opts)

	// t.Log(db.Tables[tblName])

	//t.Log(db.Tables[tblName].Rows[col1].Data)
}

func TestColumnNames(t *testing.T) {
	//db.ColumnNames("Test")
	//t.Log(db.ColumnNames("Test"))
	// t.Log(db.Tables[tblName].Rows[col1])
}

func TestUpdate(t *testing.T) {

	query := "UPDATE testTbl SET col1 = \"NEW\" WHERE col2 = \"World\";"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		// Get the Update opts with the query
		t.Log(r)
		opts := interpreter.DescribeUpdate(r)
		// test := make(map[string]string)
		// test["col1"] = "\"Hello\""
		// test["col2"] = "\"world\""
		// t.Log(opts.Condition(test))
		db.Update(opts)
	}

	t.Log(db.Tables[tblName].Rows[col1])

	// db.Update(opts)

	// t.Log(opts)

	// for k, v := range db.Tables[tblName].Rows {
	// 	t.Log(k, v)
	// }

	// t.Log(db.Tables[tblName].Rows[col1])
}

func TestDelete(t *testing.T) {

	t.Log(db.Tables[tblName].Rows["col2"])

	query := "DELETE FROM testTbl WHERE col2 = \"World\";"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		// Get the Update opts with the query
		// t.Log(r)
		opts := interpreter.DescribeDelete(r)
		// t.Log(opts)
		// test := make(map[string]string)
		// test["col1"] = "\"Hello\""
		// test["col2"] = "\"world\""
		// t.Log(opts.Condition(test))
		db.Delete(opts)
	}

	t.Log(db.Tables[tblName].Rows["col2"])
}

// func TestSelect(t *testing.T) {
// 	// query := "SELECT * FROM \"Test Table\";"
// }

func createDB() error {
	var err error
	if db, err = New(dir); err != nil {
		return err
	}

	return nil
}
