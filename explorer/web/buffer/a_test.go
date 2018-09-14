package buffer

import (
	"fmt"
	"sync"
	"testing"
)

func TestA(*testing.T) {
	a := sync.Map{}

	for i := 0; i < 100000; i++ {
		a.Store(i, i)
	}

	wg := sync.WaitGroup{}

	cnt := 0
	b := func(key, val interface{}) bool {
		cnt++
		a.Delete(key)
		_ = val
		return true
	}

	cntb := 0
	c := func(key, val interface{}) bool {
		cntb++
		return true
	}

	wg.Add(2)
	go func() {
		a.Range(b)
		fmt.Printf("cnt:%v\n\n", cnt)
		cnt = 0
		wg.Done()
	}()

	go func() {
		a.Range(c)
		fmt.Printf("cntb:%v\n\n", cntb)
		cntb = 0
		wg.Done()
	}()

	cnt = 0

	wg.Wait()
	a.Range(c)
	fmt.Printf("cnt:%v\n", cntb)
}

func TestRedis(*testing.T) {
	bb := getBlockBuffer()
	bbb := bb.GetBlock(0)
	fmt.Println(bbb)

	fmt.Println(_redisCli.Get("asdfasdf"))

}
