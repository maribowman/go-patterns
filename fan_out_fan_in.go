package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const workerCount = 3

func main() {
	inputChannel := make(chan int)
	doneChannel := make(chan struct{})

	go func() {
		time.Sleep(750 * time.Millisecond)
		// kill process
		//doneChannel <- struct{}{}
		close(doneChannel)
		fmt.Println("closed doneChannel")
	}()

	go func() {
		defer close(inputChannel)

		for i := 0; i < math.MaxInt; i++ {
			select {
			case <-doneChannel: // trigger done
				fmt.Println("killing main input")
				time.Sleep(50 * time.Millisecond)
				return
			case inputChannel <- i:
			}
		}
	}()

	workerChannel := fanOut(doneChannel, inputChannel)
	mergeChannel := fanIn(doneChannel, workerChannel...)

	for result := range randomOperation(doneChannel, mergeChannel) {
		if result%100000 == 0 {
			fmt.Println(fmt.Sprintf("run %d", result))
		}
	}
}

func fanOut(doneChannel chan struct{}, inputChannel chan int) []chan int {
	resultChannels := make([]chan int, workerCount)

	for i := 0; i < workerCount; i++ {
		workerChannel := make(chan int)
		resultChannels[i] = workerChannel

		go func(workerId int) {
			defer close(workerChannel)

			for data := range inputChannel {
				result := data + 1 // actual operation
				select {
				case <-doneChannel:
					fmt.Println(fmt.Sprintf("fan-out: shutting down worker#%d", workerId))
					time.Sleep(50 * time.Millisecond)
					return
				case workerChannel <- result:
				}
			}
		}(i)
	}

	return resultChannels
}

func fanIn(doneChannel chan struct{}, workerChannel ...chan int) chan int {
	mergeChannel := make(chan int)
	var wg sync.WaitGroup

	for channelId, channel := range workerChannel {
		wg.Add(1)
		channelClosure := channel

		go func(channelId int) {
			defer wg.Done()
			for data := range channelClosure {
				select {
				case <-doneChannel: // shutting down fan-in
					fmt.Println(fmt.Sprintf("fan-in: shutting down workerChannel#%d", channelId))
					time.Sleep(50 * time.Millisecond)
					return
				case mergeChannel <- data:
				}
			}
		}(channelId)
	}

	go func() {
		wg.Wait()
		close(mergeChannel)
	}()

	return mergeChannel
}

func randomOperation(doneChannel chan struct{}, inputChannel chan int) chan int {
	resultChannel := make(chan int)

	go func() {
		defer close(resultChannel)

		for data := range inputChannel {
			result := data - 1 // actual operation
			select {
			case <-doneChannel: // exiting random operation
				fmt.Println("shutting operation")
				time.Sleep(50 * time.Millisecond)
				return
			case resultChannel <- result:
			}
		}
	}()

	return resultChannel
}
