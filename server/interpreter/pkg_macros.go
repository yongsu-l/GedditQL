package interpreter

// const
///////////////////////////////////////////////////////////////////////////////////////////////
const (
	SELECT   = "select"
	DISTINCT = "distinct"
	ALL      = "*"
	AS       = "as"
	FROM     = "from"
	WHERE    = "where"
	ORDER    = "order"
	BY       = "by"
	LIMIT    = "limit"
	AND      = "and"
	OR       = "or"
	ASC      = "asc"
	DESC     = "desc"

	CREATE     = "create"
	TABLE      = "table"
	VIEW       = "view"
	IF         = "if"
	NOT        = "not"
	EXISTS     = "exists"
	NULL       = "null"
	PRIMARY    = "primary"
	KEY        = "key"
	UNIQUE     = "unique"
	DEFAULT    = "default"
	REFERENCES = "references"

	INSERT = "insert"
	INTO   = "into"
	VALUES = "values"

	UPDATE = "update"
	SET    = "set"

	DELETE = "delete"

	EOQ   = ";"
	COMMA = ","
	DOT   = "."

	TRUE  = "true"
	FALSE = "false"

	LPAREN    = "("
	RPAREN    = ")"
	MUL       = "*"
	DIV       = "/"
	ADD       = "+"
	SUB       = "-"
	MOD       = "%"
	EQ        = "="
	NEQ       = "!="
	LESS      = "<"
	GREATER   = ">"
	ISNOT     = "<>"
	LESSEQ    = "<="
	GREATEREQ = ">="

	ErrNoQuery    = "No query operation specified"
	ErrNoSelExp   = "Missing select expression(s)"
	ErrNoSetExp   = "Missing set expression(s)"
	ErrNoCreate   = "Cannot create"
	ErrNoColRef   = "Cannot find column reference"
	ErrNoColDef   = "Missing column definitions"
	ErrIllEOQuery = "Illegal end of query"
	ErrIllColRef  = "Illegal column reference, reserved name"
	ErrIllTabRef  = "Illegal table reference, reserved name"
	ErrIllAlias   = "Illegal alias, reserved name"
	ErrIllValue   = "Illegal value"
	ErrIllDType   = "Illegal datatype"
	ErrInvalQuery = "Invalid query operation"
	ErrUnexpToken = "Unexpected token"
	ErrExpToken   = "Expect token "
	ErrExpInt     = "Expect an integer"
	ErrOpenQuote  = "Open quotations"
	ErrValMatch   = "Type does not match value"
	ErrLenMatch   = "The number of columns does not match the number of values"
)

// reserved words
var (
	RESERVED = map[string]bool{
		"select":   true,
		"distinct": true,
		"*":        true,
		"as":       true,
		"from":     true,
		"where":    true,
		"order":    true,
		"by":       true,
		"limit":    true,
		"and":      true,
		"or":       true,
		"asc":      true,
		"desc":     true,
		"create":   true,
		"table":    true,
		"view":     true,
		"if":       true,
		"not":      true,
		"exists":   true,
		"insert":   true,
		"into":     true,
		"values":   true,
		"update":   true,
		"set":      true,
		"delete":   true,
		";":        true,
		": true,":  true,
		".":        true,
		"true":     true,
		"false":    true,
		"(":        true,
		")":        true,
		"/":        true,
		"+":        true,
		"-":        true,
		"%":        true,
		"!":        true,
		"@":        true,
		"#":        true,
		"$":        true,
		"^":        true,
		"&":        true,
		"_":        true,
		"=":        true,
		"<":        true,
		"<>":       true,
		"<=":       true,
		">":        true,
		">=":       true,
		"!=":       true,
		"{":        true,
		"}":        true,
		"|":        true,
		"\\":       true,
		"\"":       true,
		"'":        true,
		":":        true,
		"?":        true,
		"~":        true,
		"`":        true,
	}
	DATATYPES = map[string]bool{
		"float":  true,
		"int":    true,
		"bool":   true,
		"string": true,
	}
)
