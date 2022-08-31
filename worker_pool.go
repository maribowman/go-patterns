package main

import (
	"fmt"
	"sync"
)

const workerPoolSize = 3
const jobsTotal = 100

type jobDetails struct {
	jobId int
}

func main() {
	var jobs []jobDetails
	for i := 0; i < jobsTotal; i++ {
		jobs = append(jobs, jobDetails{jobId: i})
	}

	jobChannel := make(chan jobDetails)
	var wg sync.WaitGroup
	var resultMap sync.Map

	for id := 0; id < workerPoolSize; id++ {
		wg.Add(1)
		go worker(id, jobChannel, &resultMap, &wg)
	}

	for _, job := range jobs {
		jobChannel <- job
	}

	close(jobChannel)
	wg.Wait()

	var worker0, worker1, worker2 int
	resultMap.Range(func(key, value any) bool {
		switch value {
		case 0:
			worker0++
		case 1:
			worker1++
		case 2:
			worker2++
		}
		return true
	})
	fmt.Println(fmt.Sprintf("worker 0 processed %d jobs", worker0))
	fmt.Println(fmt.Sprintf("worker 1 processed %d jobs", worker1))
	fmt.Println(fmt.Sprintf("worker 2 processed %d jobs", worker2))
}

func worker(workerId int, jobChannel <-chan jobDetails, result *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobChannel {
		// actual operation
		result.Store(job.jobId, workerId)
	}
}
