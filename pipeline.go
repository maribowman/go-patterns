package main

import "fmt"

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

	var index int
	for result := range stepTwo(doneChannel, stepOne(doneChannel, inputChannel)) {
		index++
		fmt.Println(fmt.Sprintf("run %d: %d", index, result))
	}
}

func stepOne(doneChannel chan struct{}, inputChannel chan int) chan int {
	outputChannel := make(chan int)

	go func() {
		defer close(outputChannel)

		for data := range inputChannel {
			result := data + 1 // actual operation
			select {
			case <-doneChannel: // shutting down step one
				return
			case outputChannel <- result:
			}
		}
	}()

	return outputChannel
}

func stepTwo(doneChannel chan struct{}, inputChannel chan int) chan int {
	outputChannel := make(chan int)

	go func() {
		defer close(outputChannel)

		for data := range inputChannel {
			result := data - 1 // actual operation
			select {
			case <-doneChannel: // shutting down step two
				return
			case outputChannel <- result:
			}
		}
	}()

	return outputChannel
}
