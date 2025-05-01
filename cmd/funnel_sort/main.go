package main

import (
	"cse549t-project/funnelsort"
	"cse549t-project/matrix"
	"encoding/json"
	"github.com/klauspost/cpuid/v2"
	"log"
	"os"
	"runtime"
	"sort"
	"time"
)

const (
	n        = 100
	filename = "funnelsort_results"
)

type Result struct {
	Algorithm string `json:"algorithm"`
	Time      int64  `json:"time"`
	Size      int    `json:"size"`
}

func measurePerformance(f func(a []int) []int, d int) []time.Duration {
	durations := make([]time.Duration, n)
	for i := range n {
		arr := matrix.MagicMatrix1D(d, int(time.Now().UnixNano()+int64(i*7)))
		start := time.Now()
		f(arr)
		elapsed := time.Since(start)

		durations[i] = elapsed
	}
	return durations
}

func getCacheInfo() (l1, l2, l3, lineSize int) {
	l1 = cpuid.CPU.Cache.L1D
	l2 = cpuid.CPU.Cache.L2
	l3 = cpuid.CPU.Cache.L3
	lineSize = cpuid.CPU.CacheLine
	return
}

func storeResults(filename string, data []Result) error {
	// Open (or create) file for writing, overwriting if it exists
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	// Marshal JSON and write to file
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write JSON data to file
	if _, err := file.Write(bytes); err != nil {
		return err
	}

	// Close the file
	return file.Close()
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

func warmup() {
	measurePerformance(func(a []int) []int {
		c := make([]int, len(a))
		copy(c, a)
		sort.Ints(c)
		return c
	}, 50000)
}

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	fileExists := fileExists("./data/" + filename + ".json")
	if !fileExists {
		log.Println("File does not exist.")
	} else {
		log.Println("File exists.")
	}

	var results []Result
	l1, l2, l3, line := getCacheInfo()

	log.Println("Number of Processors:", runtime.GOMAXPROCS(runtime.NumCPU()))
	log.Println("L1 Cache Size:", l1)
	log.Println("L2 Cache Size:", l2)
	log.Println("L3 Cache Size:", l3)
	log.Println("Line Size:", line)

	// Warmup cache
	warmup()

	for k := 8; k <= 24; k++ {
		log.Printf("Computing array size 2^%d", k)

		s := 1 << k

		builtIn := measurePerformance(func(a []int) []int {
			sort.Ints(a)
			return a
		}, s)

		var baseline int64
		for _, duration := range builtIn {
			baseline += duration.Nanoseconds()
		}
		baseline /= int64(n)

		log.Printf("Baseline: %dms", baseline/1000/1000)

		// Append baseline result
		results = append(results, Result{
			Algorithm: "sort.Ints",
			Time:      baseline,
			Size:      s,
		})

		runTimes := measurePerformance(funnelsort.Sort, s)

		var avg int64
		// Append each run time result
		for _, runTime := range runTimes {
			avg += runTime.Nanoseconds()
			results = append(results, Result{
				Algorithm: "FunnelSort",
				Time:      runTime.Nanoseconds(),
				Size:      s,
			})
		}
		avg /= int64(n)

		log.Printf("FunnelSort speedup: %f", float64(baseline)/float64(avg))
	}

	// -----------------------------------------------------
	// Store Results
	// -----------------------------------------------------
	err := storeResults("./data/"+filename+".json", results)
	if err != nil {
		log.Println("❌ Failed to write:", err)
	} else {
		log.Println("✅ Results written.")
	}
}
