package mysql

import "testing"

func TestDistinct(t *testing.T) {
	var data = []string{"a", "ab", "ac", "a1", "a2", "a", "a"}
	ret, distinct := Distinct(data)
	t.Log(ret, "\n")
	t.Log(distinct, "\n")
}
