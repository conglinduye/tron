package buffer

import (
	"sync"
)

type dbIndex struct {
	sync.RWMutex
	index []*indexPos
	total int64
	step  int64
}

type indexPos struct {
	Position int64
	BlockID  int64
	Offset   int64
	Count    int64
}

// GetTotal ...
func (dbi *dbIndex) GetTotal() int64 {
	dbi.RLock()
	defer dbi.RUnlock()
	return dbi.total
}

// GetStep ...
func (dbi *dbIndex) GetStep() int64 {
	dbi.RLock()
	defer dbi.RUnlock()
	return dbi.step
}

// GetIndex ...
func (dbi *dbIndex) GetIndex() []*indexPos {
	dbi.RLock()
	defer dbi.RUnlock()
	return dbi.index
}
