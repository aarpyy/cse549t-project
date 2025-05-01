package main

import (
	"cse549t-project/matmul"
	"cse549t-project/matrix"
	"encoding/json"
	"fmt"
	"github.com/klauspost/cpuid/v2"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	n    = 100
	dim1 = 1024
	dim2 = 2000
	dim3 = 1024
)

type Result struct {
	Algorithm        string `json:"algorithm"`
	Time             int64  `json:"time"`
	AssumedCacheSize int    `json:"assumed_cache_size"`
}

func measurePerformance(f func(a, b [][]int) [][]int) []time.Duration {
	durations := make([]time.Duration, n)
	for i := range n {
		a := matrix.MagicMatrix(dim1, dim2, int(time.Now().UnixNano()+int64(i*7)))
		b := matrix.MagicMatrix(dim2, dim3, int(time.Now().UnixNano()+int64(i*7)))
		start := time.Now()
		f(a, b)
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
	measurePerformance(matmul.MatrixMultiply)
}

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	filename := fmt.Sprintf("cache_aware_mm_%dx%dx%d", dim1, dim2, dim3)

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

	// Also compare for slightly smaller than cache size to account for randomness
	cacheConfigs := []string{
		"4KB",
		"8KB",
		"16KB",
		"32KB",
		"64KB",
		"128KB",
		"256KB",
		"512KB",
		"1MB",
		"2MB",
		"4MB",
		"8MB",
		"16MB",
		"32MB",
		"64MB",
		"128MB",
	}

	// Compute average
	simple := measurePerformance(matmul.MatrixMultiply)
	var baseline int64
	for _, duration := range simple {
		baseline += duration.Nanoseconds()
	}
	baseline /= int64(n)
	log.Printf("Baseline: %dms", baseline/1000/1000)

	// Append baseline result
	results = append(results, Result{
		Algorithm:        "TripleNestedMM",
		Time:             baseline,
		AssumedCacheSize: -1,
	})

	for _, cfg := range cacheConfigs {
		log.Printf("Assuming cache size: %s", cfg)

		B := matmul.ComputeBlockSize(cfg)
		runTimes := measurePerformance(func(a, b [][]int) [][]int {
			return matmul.CacheAwareMM(a, b, B)
		})

		var avg int64
		// Append each run time result
		for _, runTime := range runTimes {
			avg += runTime.Nanoseconds()
			results = append(results, Result{
				Algorithm:        "CacheAwareMM",
				Time:             runTime.Nanoseconds(),
				AssumedCacheSize: B,
			})
		}
		avg /= int64(n)

		log.Printf("CacheAwareMM speedup: %f", float64(baseline)/float64(avg))
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
