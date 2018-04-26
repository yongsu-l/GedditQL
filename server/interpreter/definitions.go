package interpreter

///////////////////////////////////////////////////////////////////////////////////////////////
type selectOptions struct {
	Distinct   bool
	All        bool
	ColumnRefs []string
	As         map[string]string
	TableRefs  []string
	Condition  func(map[string]string) bool
	Order      string
	Limit      int
}
