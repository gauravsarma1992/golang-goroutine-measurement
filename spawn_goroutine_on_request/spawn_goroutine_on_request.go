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
	OutputChSize = 100000
)

func calculateSum(outputCh chan int, a, b int) {
	result := 0
	for currIdx := 0; currIdx < 100; currIdx++ {
		result = a + b
	}
	outputCh <- result
	return
}

func main() {
	var (
		frequency int
	)

	outputCh := make(chan int, OutputChSize)
	flag.Parse()
	frequency, _ = strconv.Atoi(flag.Args()[0])

	returnedCount := 0

	f, err := os.Create("goroutine.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	log.Println("Frequency selected", frequency)
	startTime := time.Now()

	go func() {
		for currIdx := 0; currIdx < frequency; currIdx++ {
			go calculateSum(outputCh, 4, currIdx)
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
