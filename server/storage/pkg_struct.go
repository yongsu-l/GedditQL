package storage

// Response is the type which is first an array of column names, then an array of types, then an array of all the data
type Response struct {
	Names        []string
	DataTypes    []string
	Data         [][]string
	RowsAffected int
	Result       string
	Err          string
}

// Data which will hold all of the data
type Data struct {
	DataType string
	Data     []string
}

// Table which holds all of the values in an array as order matters
type Table struct {
	// Data held within each row linearly
	// Key is column name
	Rows   map[string]*Data
	Length int
}

// Database which holds all of the tables
type Database struct {
	Dir    string
	File   string
	Tables map[string]*Table
}
