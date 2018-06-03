package sorter

import (
	"sort"
)

// Stat is a data structure to hold key/value pairs
// method taken from https://github.com/indraniel/go-learn/blob/master/09-sort-map-keys-by-values.go
type Stat struct {
	Key   string
	Value int
}

// StatList is slice of pairs that implements sort.Interface to sort by values
type StatList []Stat

func (st StatList) Len() int           { return len(st) }
func (st StatList) Swap(i, j int)      { st[i], st[j] = st[j], st[i] }
func (st StatList) Less(i, j int) bool { return st[i].Value < st[j].Value }

// Sorter provides a descending slice of key value pairs from a map
func Sorter(input map[string]int) StatList {
	stats := make(StatList, len(input))
	i := 0
	for k, v := range input {
		stats[i] = Stat{k, v}
		i++
	}

	sort.Sort(sort.Reverse(stats))

	return stats
}
