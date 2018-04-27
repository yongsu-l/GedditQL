package interpreter

import "strings"

func getColumnDef(tq *queue) *ColumnDef {
	name := tq.Current()
	dType := strings.ToLower(tq.Next().Current())
	tq.Next()
	notNull := false
	primaryKey := false
	unique := false
	defalt := ""

	for tq.Current() != COMMA && tq.Current() != RPAREN {
		switch tk := strings.ToLower(tq.Current()); tk {
		case NOT:
			notNull = true
			tq.Next().Next()
			break
		case PRIMARY:
			primaryKey = true
			tq.Next().Next()
			break
		case UNIQUE:
			unique = true
			tq.Next()
			break
		case DEFAULT:
			_, val, _ := getValue(tq.Next())
			defalt = val
			break
		default:
			break
		}
	}

	return &ColumnDef{
		Name:       name,
		DataType:   dType,
		NotNull:    notNull,
		PrimaryKey: primaryKey,
		Unique:     unique,
		Default:    defalt,
	}
}
