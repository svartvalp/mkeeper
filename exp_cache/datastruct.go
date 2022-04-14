package exp_cache

import (
	"container/list"
)

type entry struct {
	Key   interface{}
	Value interface{}
	Hash  uint64
	Exp   int64
	ExpEl *list.Element
}
