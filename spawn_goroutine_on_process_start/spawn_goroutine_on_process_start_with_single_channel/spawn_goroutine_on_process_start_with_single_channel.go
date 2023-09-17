package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

const (
	ConcurrencyCount = 10
)

var (
	inputCh  chan int
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

func processRequest() {
	for {
		randInt := <-inputCh
		calculateSum(4, randInt)
	}
	return
}

func startGoroutines(concurrencyCount int) {
	for currIdx := 0; currIdx < concurrencyCount; currIdx++ {
		go processRequest()
	}
}

func main() {
	var (
		frequency int
	)

	inputCh = make(chan int, 1024)
	outputCh = make(chan int, 1024)

	flag.Parse()
	frequency, _ = strconv.Atoi(flag.Args()[0])

	returnedCount := 0

	startGoroutines(ConcurrencyCount)
	time.Sleep(5 * time.Second)

	log.Println("Frequency selected", frequency)

	f, err := os.Create("goroutine.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	startTime := time.Now()

	go func() {
		for currIdx := 0; currIdx < frequency; currIdx++ {
			inputCh <- currIdx
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
