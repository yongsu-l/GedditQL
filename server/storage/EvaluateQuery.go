package storage

import (
	"GedditQL/server/interpreter"
	"strings"
)

//EvaluateQuery is a switch statement which checks the token and calls on the appropriate Describe functions
func (db *Database) EvaluateQuery(tks []string) (Response, error) {

	switch tk := strings.ToLower(tks[0]); tk {
	case SELECT:
		return db.Select(interpreter.DescribeSelect(tks))
	case UPDATE:
		return db.Update(interpreter.DescribeUpdate(tks))
	case DELETE:
		return db.Delete(interpreter.DescribeDelete(tks))
	case INSERT:
		return db.Insert(interpreter.DescribeInsert(tks))
	case CREATE:
		return db.Create(interpreter.DescribeCreate(tks))
	default:
		return Response{}, &errorString{"Unsupported query"}
	}
}
