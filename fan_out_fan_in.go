package main

import (
	"fmt"
	"sync"
)

const workerCount = 5

func main() {
	input := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	inputChannel := make(chan int)
	doneChannel := make(chan struct{})
	defer close(doneChannel)

	go func() {
		defer close(inputChannel)

		for _, data := range input {
			select {
			case <-doneChannel: // trigger done
				return
			case inputChannel <- data:
			}
		}
	}()

	channels := fanOut(doneChannel, inputChannel)
	mergeChannel := fanIn(doneChannel, channels...)

	var index int
	for randomNumber := range randomOperation(doneChannel, mergeChannel) {
		index++
		fmt.Printf("run %d: %d\n", index, randomNumber)
	}
}

func fanOut(doneChannel chan struct{}, inputChannel chan int) []chan int {
	resultChannels := make([]chan int, workerCount)

	for i := 0; i < workerCount; i++ {
		workerChannel := make(chan int)
		resultChannels[i] = workerChannel

		go func() {
			defer close(workerChannel)
			for data := range inputChannel {
				result := data + 1 // actual operation
				select {
				case <-doneChannel: // shutting down fan-out
					return
				case workerChannel <- result:
				}
			}
		}()
	}

	return resultChannels
}

func fanIn(doneChannel chan struct{}, resultChannels ...chan int) chan int {
	mergeChannel := make(chan int)
	var wg sync.WaitGroup

	for _, channel := range resultChannels {
		wg.Add(1)
		channelClosure := channel

		go func() {
			defer wg.Done()
			for data := range channelClosure {
				select {
				case <-doneChannel: // shutting down fan-in
					return
				case mergeChannel <- data:
				}
			}
		}()
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
				return
			case resultChannel <- result:
			}
		}
	}()

	return resultChannel
}
