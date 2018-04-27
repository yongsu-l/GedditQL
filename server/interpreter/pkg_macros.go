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
	ErrInvalQuery = "Invalid query operation"
	ErrIllEOQuery = "Illegal end of query"
	ErrExpToken   = "Expect token "
	ErrUnexpToken = "Unexpected token"
	ErrExpInt     = "Expect an integer"
	ErrNoSelExp   = "Missing select expression(s)"
	ErrNoSetExp   = "Missing set expression(s)"
	ErrOpenQuote  = "Open quotations"
	ErrIllColRef  = "Illegal column reference, reserved name"
	ErrIllTabRef  = "Illegal table reference, reserved name"
	ErrIllAlias   = "Illegal alias, reserved name"
	ErrNoCreate   = "Cannot create"
	ErrIllValue   = "Illegal value"
	ErrIllDType   = "Illegal datatype"
	ErrNoColRef   = "Cannot find column reference"
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
