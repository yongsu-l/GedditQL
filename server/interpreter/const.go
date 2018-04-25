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
)

// reserved words
var (
	RESERVED = []string{
		"select",
		"distinct",
		"*",
		"as",
		"from",
		"where",
		"order",
		"by",
		"limit",
		"and",
		"or",
		"asc",
		"desc",
		"create",
		"table",
		"view",
		"if",
		"not",
		"exists",
		"insert",
		"into",
		"values",
		"update",
		"set",
		"delete",
		";",
		",",
		".",
		"true",
		"false",
		"(",
		")",
		"/",
		"+",
		"-",
		"%",
		"!",
		"@",
		"#",
		"$",
		"^",
		"&",
		"_",
		"=",
		"<",
		"<>",
		"<=",
		">",
		">=",
		"!=",
		"{",
		"}",
		"|",
		"\\",
		"\"",
		"'",
		":",
		"/",
		"?",
		"~",
		"`",
	}
	DATATYPES = []string{"float", "int", "bool", "string"}
)
