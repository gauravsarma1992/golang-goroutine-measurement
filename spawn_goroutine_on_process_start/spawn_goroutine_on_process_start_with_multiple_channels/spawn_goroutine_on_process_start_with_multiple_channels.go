package main

import (
	"flag"
	"log"
	"strconv"
	"time"
)

const (
	ConcurrencyCount = 10000
)

var (
	inputChs []chan int
	outputCh chan int
)

func calculateSum(a, b int) {
	result := 0
	for currIdx := 0; currIdx < 100; currIdx++ {
		result = a + b
	}
	outputCh <- result
	return
}

func processRequest(inputCh chan int) {
	for {
		randInt := <-inputCh
		calculateSum(4, randInt)
	}
	return
}

func startGoroutines(concurrencyCount int) {
	for currIdx := 0; currIdx < concurrencyCount; currIdx++ {
		inputCh := make(chan int, 1024)
		inputChs[currIdx] = inputCh
		go processRequest(inputCh)
	}
}

func main() {
	var (
		frequency int
	)

	inputChs = make([]chan int, ConcurrencyCount)
	outputCh = make(chan int, 1024)

	flag.Parse()
	frequency, _ = strconv.Atoi(flag.Args()[0])

	returnedCount := 0
	chanIdx := 0

	startGoroutines(ConcurrencyCount)
	time.Sleep(5 * time.Second)

	log.Println("Frequency selected", frequency)
	startTime := time.Now()

	go func() {
		for currIdx := 0; currIdx < frequency; currIdx++ {
			if chanIdx >= ConcurrencyCount {
				chanIdx = 0
			}
			inputChs[chanIdx] <- currIdx
			chanIdx += 1
		}
	}()

	for {
		<-outputCh
		returnedCount += 1
		if returnedCount == frequency {
			break
		}
	}
	timeTaken := time.Since(startTime)
	log.Println("Total Queries processed:", returnedCount, "in", timeTaken.Milliseconds(), "ms")

}
