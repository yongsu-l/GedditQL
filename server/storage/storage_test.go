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

func TestCreate(t *testing.T) {
	query := "CREATE table test2 (col1 string, col2 string);"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		opts := interpreter.DescribeCreate(r)
		db.Create(opts)
		if db.Tables["test2"] == nil {
			t.Fatal("Failed to create table")
		}
	}
}

func TestInsert(t *testing.T) {
	query := "INSERT INTO testTbl (col1, col2) VALUES (\"hello\", \"world\");"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		prevLength := db.Tables["testTbl"].Length
		opts := interpreter.DescribeInsert(r)
		db.Insert(opts)
		if prevLength+1 != db.Tables["testTbl"].Length {
			t.Fatal("Failed to insert")
		}
	}

	query = "INSERT INTO testTbl (col1, col2) VALUES (\"anakin\", \"world\");"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		prevLength := db.Tables["testTbl"].Length
		opts := interpreter.DescribeInsert(r)
		db.Insert(opts)
		if prevLength+1 != db.Tables["testTbl"].Length {
			t.Fatal("Failed to insert")
		}
	}

	query = "INSERT INTO testTbl (col1, col2) VALUES (\"skywalker\", \"world\");"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		prevLength := db.Tables["testTbl"].Length
		opts := interpreter.DescribeInsert(r)
		db.Insert(opts)
		if prevLength+1 != db.Tables["testTbl"].Length {
			t.Fatal("Failed to insert")
		}
	}
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
		// t.Log(r)
		opts := interpreter.DescribeUpdate(r)
		// test := make(map[string]string)
		// test["col1"] = "\"Hello\""
		// test["col2"] = "\"world\""
		// t.Log(opts.Condition(test))
		db.Update(opts)
	}

	// t.Log(db.Tables[tblName].Rows[col1])

	// db.Update(opts)

	// t.Log(opts)

	// for k, v := range db.Tables[tblName].Rows {
	// 	t.Log(k, v)
	// }

	// t.Log(db.Tables[tblName].Rows[col1])
}

func TestSelect(t *testing.T) {

	query := "SELECT * FROM testTbl;"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		opts := interpreter.DescribeSelect(r)
		res, err := db.Select(opts)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log(res)
		}
	}

	query = "SELECT col1 FROM testTbl where col1 =\"hello\";"
	t.Log(query)
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		// t.Log(r)
		opts := interpreter.DescribeSelect(r)
		res, err := db.Select(opts)
		t.Log(res, err)
	}

	//Second query where there will be an error
	query = "SELECT ERROR FROM testTbl where col1 =\"hello\";"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		opts := interpreter.DescribeSelect(r)
		_, err := db.Select(opts)
		if err == nil {
			t.Fatal("Should be expected error")
		}
	}
	query = "SELECT 1;"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		opts := interpreter.DescribeSelect(r)
		res, _ := db.Select(opts)
		t.Log(res)
	}

	query = "SELECT * FROM testTbl;"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		opts := interpreter.DescribeSelect(r)
		res, _ := db.Select(opts)
		t.Log(res)
	}
}

func TestSelectLimit(t *testing.T) {
	// Check if limit actually limits the data
	query := "SELECT * FROM testTbl LIMIT 1;"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		opts := interpreter.DescribeSelect(r)
		if res, err := db.Select(opts); err != nil {
			t.Fatal(err)
		} else if len(res.Data) != 1 || res.RowsAffected != 1 {
			t.Fatal("Failed to limit")
		}
	}
}

func TestSelectFunction(t *testing.T) {
	// First test sum of columns
	query := "SELECT SUM(col1) from testTbl LIMIT 1;"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		opts := interpreter.DescribeSelect(r)
		if _, err := db.Select(opts); err == nil {
			t.Fatal("Should return err because summing strings")
		}
	}

	query = "SELECT COUNT(col1) FROM testTbl;"
	if r, err := parser.Tokenize(query); err != nil {
		t.Fatal(err)
	} else {
		opts := interpreter.DescribeSelect(r)
		if res, err := db.Select(opts); err != nil {
			t.Fatal(err)
		} else {
			t.Log(res)
		}
	}
}

func createDB() error {
	var err error
	if db, err = New(dir); err != nil {
		return err
	}

	return nil
}
