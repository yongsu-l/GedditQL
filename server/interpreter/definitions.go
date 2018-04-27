package interpreter

///////////////////////////////////////////////////////////////////////////////////////////////
type SelectOptions struct {
	Distinct   bool
	All        bool
	ColumnRefs []string
	As         map[string]string
	TableRefs  []string
	Condition  func(map[string]string) (bool, error)
	Order      string
	By         string
	Limit      int
}
