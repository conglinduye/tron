package service

import (
	"fmt"
	"testing"
)

func TestGetNextCycle(t *testing.T) {
	tt, _ := QueryVoteNextCycle()
	fmt.Sprintf("%v", tt.NextCycle)
}

func TestGetMarket(t *testing.T) {
	/*	ss, err := QueryMarkets()
		fmt.Sprintf("%v,%v", ss, err)*/
}
