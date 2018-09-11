package service

import (
	"fmt"
	"testing"
)

func TestGetNextCycle(t *testing.T) {
	tt, _ := QueryVoteNextCycle()
	fmt.Sprintf("%v", tt.NextCycle)
}
