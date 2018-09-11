package main

type workerCounter struct {
	num int
	c   chan struct{}
}

func (wc *workerCounter) startOne() {
	<-wc.c
}

func (wc *workerCounter) stopOne() {
	wc.c <- struct{}{}
}

func (wc *workerCounter) currentWorker() int {
	return wc.num - len(wc.c)
}

// newWorkerCounter ...
func newWorkerCounter(maxWorker int) *workerCounter {
	ret := &workerCounter{
		num: maxWorker,
		c:   make(chan struct{}, maxWorker),
	}

	for i := 0; i < ret.num; i++ {
		ret.stopOne()
	}

	return ret
}
