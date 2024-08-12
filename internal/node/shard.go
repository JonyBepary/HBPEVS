package node

type LookupTable struct {
	table map[string]int
}

func NewLookupTable() *LookupTable {
	return &LookupTable{
		table: make(map[string]int),
	}
}

func (lt *LookupTable) Add(key string, value int) {
	lt.table[key] = value
}

func (lt *LookupTable) Get(key string) int {
	value, ok := lt.table[key]
	if !ok {
		return -1
	}

	return value
}
