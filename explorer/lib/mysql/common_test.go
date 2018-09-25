package mysql

import (
	"fmt"
	"testing"
)

func TestDistinct(t *testing.T) {
	var data = []string{"a", "ab", "ac", "a1", "a2", "a", "a"}
	ret, distinct := Distinct(data)
	t.Log(ret, "\n")
	t.Log(distinct, "\n")
}

func TestRemoveAddr(t *testing.T) {
	ss := []string{"11", "bb", "fcc"}
	for key, val := range ss {
		if val == "fcc" {
			kk := key + 1
			ss = append(ss[:key], ss[kk:]...)
		}
	}
	fmt.Printf("ss:[%v]", ss)
}
