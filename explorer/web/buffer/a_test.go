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
		log.Infof("cnt:%v\n\n", cnt)
		cnt = 0
		wg.Done()
	}()

	go func() {
		a.Range(c)
		log.Infof("cntb:%v\n\n", cntb)
		cntb = 0
		wg.Done()
	}()

	cnt = 0

	wg.Wait()
	a.Range(c)
	log.Infof("cnt:%v\n", cntb)
}

func TestRedis(*testing.T) {
	// // bb := getBlockBuffer()
	// bbb := bb.GetBlock(0)
	// fmt.Println(bbb)

	// fmt.Println(_redisCli.Get("asdfasdf"))

	a := make([]int, 10, 10)
	b := make([]int, 100, 100)

	for i := 0; i < 100; i++ {
		b[i] = i
	}

	copy(a, b[0:30])
	fmt.Println(a)

}
