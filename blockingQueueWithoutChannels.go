package main

import (
	"sync"
	"time"
)

type BlockingQueue struct {
	capacity int
	mu       sync.Mutex
	cond     sync.Cond
	queue    []int
}

func NewBlockingQueue(capacity int) *BlockingQueue {
	blockingQueue := &BlockingQueue{}
	blockingQueue.capacity = capacity
	blockingQueue.queue = make([]int, 0, capacity)
	blockingQueue.cond = sync.Cond{L: &blockingQueue.mu}
	return blockingQueue
}

func (bq *BlockingQueue) Enqueue(item int) {
	bq.mu.Lock()
	defer bq.mu.Unlock()
	for len(bq.queue) == bq.capacity {
		bq.cond.Wait()
	}
	bq.queue = append(bq.queue, item)
	bq.cond.Signal()
}

func (bq *BlockingQueue) Dequeue() int {
	bq.mu.Lock()
	defer bq.mu.Unlock()
	for len(bq.queue) == 0 {
		bq.cond.Wait()
	}
	item := bq.queue[0]
	bq.queue = bq.queue[1:]
	bq.cond.Signal()
	return item
}

func main() {
	bq := NewBlockingQueue(10)

	go func() {
		for i := 0; i < 100; i++ {
			//time.Sleep(1 * time.Millisecond)
			bq.Enqueue(i)
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			//time.Sleep(10 * time.Millisecond)
			println(bq.Dequeue())
		}
	}()
	time.Sleep(1 * time.Minute)
}
