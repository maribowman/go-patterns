package main

import (
	"fmt"
	"sync"
	"time"
)

const requestCount = 5

type Semaphore struct {
	semaCh chan struct{}
}

func NewSemaphore(maxRequests int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, maxRequests),
	}
}

func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.semaCh
}

func main() {
	var wg sync.WaitGroup
	semaphore := NewSemaphore(requestCount)

	for index := 0; index < 15; index++ {
		wg.Add(1)

		go func(taskID int) {
			semaphore.Acquire()
			defer wg.Done()
			defer semaphore.Release()

			fmt.Println(fmt.Sprintf("%s Running worker %d", time.Now().Format("15:04:05"), taskID))
			time.Sleep(1 * time.Second)
		}(index)
	}

	wg.Wait()
}
