package storage

import (
	"GedditQL/server/interpreter"
	//"log"
	"os"
	"testing"
)

var (
	db      *Database
	dir     = "test"
	tblName = "Test Table"
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
	data := &Data{Type: "string"}
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

func TestColumnNames(t *testing.T) {
	db.ColumnNames("Test")
}

func TestInsertInto(t *testing.T) {
	insertion := make(map[string]string)

	insertion[col1] = "hello"
	insertion[col2] = "world"

	db.InsertInto(tblName, insertion)

	insertion = make(map[string]string)
	insertion[col1] = "second round"

	db.InsertInto(tblName, insertion)
}

// func TestSelect(t *testing.T) {
//
// 	colNames := []string{col1, "Test column"}
//
// 	if _, err := db.Select(colNames, tblName); err == nil {
// 		t.Error("Should spit error because column doesn't exist")
// 	}
//
// 	colNames = []string{col1, col2}
//
// 	if tbl, err := db.Select(colNames, tblName); err != nil {
// 		t.Error("Shouldn't spit out error because columns exist")
// 	} else if len(tbl.Data) != 2 {
// 		t.Error("Should be returning two rows of data set")
// 	}
// }

func TestSelect(t *testing.T) {
	opts := interpreter.SelectOptions{}
	opts.TableRefs = append(opts.TableRefs, tblName)

	opts.All = true

	//log.Println(opts)

	if _, err := db.Select(opts); err != nil {
		t.Fatal(err)
	}

	opts.Distinct = true

	insertion := make(map[string]string)

	insertion[col1] = "hello"
	insertion[col2] = "world"

	db.InsertInto(tblName, insertion)

	if tbl, err := db.Select(opts); err != nil {
		t.Fatal("Error selecting")
	} else if len(tbl.Data) != 2 {
		t.Fatal("Error choosing distinct")
	}

	opts = interpreter.SelectOptions{
		TableRefs:  []string{tblName},
		ColumnRefs: []string{col1, col2},
	}

	if tbl, err := db.Select(opts); err != nil {
		t.Fatal(err)
	} else {
		t.Log(tbl)
	}

	// t.Log(opts.As)
	opts.As = make(map[string]string)
	opts.As[col1] = "Change Col"

	if tbl, err := db.Select(opts); err != nil {
		t.Fatal(err)
	} else if tbl.Names[0] != "Change Col" {
		t.Fatal("Failed to change column name")
	} else {
		t.Log(tbl)
	}

	t.Log(opts.Limit)

}

func createDB() error {
	var err error
	if db, err = New(dir); err != nil {
		return err
	}

	return nil
}
