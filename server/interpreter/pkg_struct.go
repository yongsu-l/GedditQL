package interpreter

// SelectOptions ...
///////////////////////////////////////////////////////////////////////////////////////////////
type SelectOptions struct {
	Distinct   bool
	All        bool
	ColumnRefs []string
	As         map[string]string
	FuncCols   []string
	FuncMap    map[string]string
	TableRefs  []string
	Condition  func(map[string]string) (bool, error)
	Order      string
	By         string
	Limit      int
}

// CreateOptions ...
///////////////////////////////////////////////////////////////////////////////////////////////
type CreateOptions struct {
	Type        string
	IfNotExists bool
	TableRef    string
	ColumnDefs  []*ColumnDef
}

// ColumnDef ...
///////////////////////////////////////////////////////////////////////////////////////////////
type ColumnDef struct {
	Name       string
	DataType   string
	NotNull    bool
	PrimaryKey bool
	Unique     bool
	Default    string
}

// InsertOptions ...
///////////////////////////////////////////////////////////////////////////////////////////////
type InsertOptions struct {
	TableRef   string
	ColumnRefs []string
	Values     []string
	ValueTypes []string
}

// UpdateOptions ...
///////////////////////////////////////////////////////////////////////////////////////////////
type UpdateOptions struct {
	TableRef  string
	ValueMap  map[string]string
	TypeMap   map[string]string
	Condition func(map[string]string) (bool, error)
}

// DeleteOptions ...
///////////////////////////////////////////////////////////////////////////////////////////////
type DeleteOptions struct {
	TableRef  string
	Condition func(map[string]string) (bool, error)
	Limit     int
}

// queue
///////////////////////////////////////////////////////////////////////////////////////////////
type queue struct {
	q []string
	p int
}

func newQueue(lst []string) *queue {
	return &queue{q: lst, p: 1}
}

func (q *queue) Next() *queue {
	q.q = q.q[1:]
	q.p++
	return q
}

func (q *queue) Current() string {
	return q.q[0]
}

func (q *queue) Pos() int {
	return q.p
}

func (q *queue) Ind(n int) string {
	return q.q[n]
}
