package main

import (
	"cse549t-project/matmul"
	"cse549t-project/matrix"
	"cse549t-project/mergesort"
	"encoding/csv"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/klauspost/cpuid/v2"
)

// Measures:
// n = 10, do this 3 times for each data set
// dimensions to test for MM = 512, 1024, 2048, 4096, 8192
// dimensions to test for MS = 4096, 8192, 16384, 32768
const (
	n        = 10
	dim1     = 4096
	dim2     = 4096
	dim3     = 4096
	filename = "olivia-pc"
)

type MatrixPair struct {
	A [][]int
	B [][]int
}

// pre-generate 10 random matrix pairs to reduce setup time
func generateMatrixPairs() []MatrixPair {
	pairs := make([]MatrixPair, n)
	for i := 0; i < n; i++ {
		pairs[i] = MatrixPair{
			A: matrix.MagicMatrix(dim1, dim2, int(time.Now().UnixNano()+int64(i*3))),
			B: matrix.MagicMatrix(dim1, dim2, int(time.Now().UnixNano()+int64(i*7))),
		}
	}
	return pairs
}

func generateArrays() [][]int {
	arrays := make([][]int, n)
	for i := 0; i < n; i++ {
		arrays[i] = matrix.MagicMatrix1D(dim1, int(time.Now().UnixNano()+int64(i*7)))
	}
	return arrays
}

func measurePerformance(f func(a, b [][]int) [][]int, pairs []MatrixPair) time.Duration {
	total := time.Duration(0)
	for _, p := range pairs {
		start := time.Now()
		f(p.A, p.B)
		elapsed := time.Since(start)
		total += elapsed
	}
	return total / time.Duration(len(pairs))
}

func measurePerformance1D(f func(a []int) []int, arrays [][]int) time.Duration {
	total := time.Duration(0)
	for _, original := range arrays {
		input := make([]int, dim1)
		copy(input, original)
		start := time.Now()
		f(input)
		elapsed := time.Since(start)
		total += elapsed
	}
	// log.Println("Total = ", total)
	// log.Println("time.Duration(n) = ", time.Duration(n))
	return total / time.Duration(n)
}

func getCacheInfo() (l1, l2, l3, lineSize int) {
	l1 = cpuid.CPU.Cache.L1D
	l2 = cpuid.CPU.Cache.L2
	l3 = cpuid.CPU.Cache.L3
	lineSize = cpuid.CPU.CacheLine
	return
}

func storeResults(filename string, data [][]string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	fileExists := fileExists("./data/" + filename + ".csv")
	if !fileExists {
		log.Println("File does not exist.")
	} else {
		log.Println("File exists.")
	}

	var results [][]string
	l1, l2, l3, line := getCacheInfo()
	numProc := runtime.GOMAXPROCS(runtime.NumCPU())
	pairs := generateMatrixPairs()
	arrays := generateArrays()

	log.Println("Number of Processors:", runtime.GOMAXPROCS(runtime.NumCPU()))
	a := matrix.MagicMatrix(dim1, dim2, 42)
	b := matrix.MagicMatrix(dim2, dim3, 21)
	log.Println("Matrix A Size: ", len(a), "x", len(a[0]))
	log.Println("Matrix B Size: ", len(b), "x", len(b[0]))

	// Nested For Loop MM
	c := matmul.MatrixMultiply(a, b)
	baseline := measurePerformance(matmul.MatrixMultiply, pairs)
	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication check failed")
	}
	log.Printf("Nested For Loop MM running time: %dμs", baseline.Microseconds())
	result := []string{
		time.Now().Format("2006-01-02 15:04:05"),
		"MM-Simple",
		strconv.Itoa(dim1),
		strconv.Itoa(dim2),
		strconv.FormatInt(baseline.Nanoseconds(), 10),
		strconv.Itoa(1),
		strconv.Itoa(numProc),
		strconv.Itoa(l1),
		strconv.Itoa(l2),
		strconv.Itoa(l3),
		strconv.Itoa(line),
		"",
		strconv.Itoa(n),
	}
	results = append(results, result)

	// Divide and Conquer MM
	c = matmul.MatrixMultiplyDandCRec(a, b, len(a))
	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer check failed")
	}
	runTime := measurePerformance(func(a, b [][]int) [][]int {
		return matmul.MatrixMultiplyDandCRec(a, b, len(a))
	}, pairs)
	log.Println("Divide and Conquer MM speedup: ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))
	result = []string{
		time.Now().Format("2006-01-02 15:04:05"),
		"MM-DivideConquer",
		strconv.Itoa(dim1),
		strconv.Itoa(dim2),
		strconv.FormatInt(runTime.Nanoseconds(), 10),
		strconv.FormatFloat(float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()), 'f', -1, 64),
		strconv.Itoa(numProc),
		strconv.Itoa(l1),
		strconv.Itoa(l2),
		strconv.Itoa(l3),
		strconv.Itoa(line),
		"",
		strconv.Itoa(n),
	}
	results = append(results, result)

	// 8x Parallel Divide and Conquer MM
	c = matmul.MatrixMultiplyDandCRecParallel8(a, b, len(a))
	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer 8x Parallel check failed")
	}
	runTime = measurePerformance(func(a, b [][]int) [][]int {
		return matmul.MatrixMultiplyDandCRecParallel8(a, b, len(a))
	}, pairs)
	log.Println("Divide and Conquer 8x Parallel MM speedup: ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))
	result = []string{
		time.Now().Format("2006-01-02 15:04:05"),
		"MM-DivideConquer-Par8",
		strconv.Itoa(dim1),
		strconv.Itoa(dim2),
		strconv.FormatInt(runTime.Nanoseconds(), 10),
		strconv.FormatFloat(float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()), 'f', -1, 64),
		strconv.Itoa(numProc),
		strconv.Itoa(l1),
		strconv.Itoa(l2),
		strconv.Itoa(l3),
		strconv.Itoa(line),
		"",
		strconv.Itoa(n),
	}
	results = append(results, result)

	// 4x Parallel Divide and Conquer MM
	c = matmul.MatrixMultiplyDandCRecParallel4(a, b, len(a))
	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer 4x Parallel check failed")
	}
	runTime = measurePerformance(func(a, b [][]int) [][]int {
		return matmul.MatrixMultiplyDandCRecParallel4(a, b, len(a))
	}, pairs)
	log.Println("Divide and Conquer 4x Parallel MM speedup: ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))
	result = []string{
		time.Now().Format("2006-01-02 15:04:05"),
		"MM-DivideConquer-Par4",
		strconv.Itoa(dim1),
		strconv.Itoa(dim2),
		strconv.FormatInt(runTime.Nanoseconds(), 10),
		strconv.FormatFloat(float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()), 'f', -1, 64),
		strconv.Itoa(numProc),
		strconv.Itoa(l1),
		strconv.Itoa(l2),
		strconv.Itoa(l3),
		strconv.Itoa(line),
		"",
		strconv.Itoa(n),
	}
	results = append(results, result)

	// Cache Aware MM
	c = matmul.CacheAwareMM(a, b, matmul.ComputeBlockSize("1MB"))
	if !matmul.CheckMatMul(a, b, c) {
		panic("Cache aware matrix multiplication check failed")
	}
	runTime = measurePerformance(func(a, b [][]int) [][]int {
		return matmul.CacheAwareMM(a, b, matmul.ComputeBlockSize("1MB"))
	}, pairs)
	log.Println("Cache aware matrix multiplication speedup (1MB): ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))
	result = []string{
		time.Now().Format("2006-01-02 15:04:05"),
		"MM-CacheAware",
		strconv.Itoa(dim1),
		strconv.Itoa(dim2),
		strconv.FormatInt(runTime.Nanoseconds(), 10),
		strconv.FormatFloat(float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()), 'f', -1, 64),
		strconv.Itoa(numProc),
		strconv.Itoa(l1),
		strconv.Itoa(l2),
		strconv.Itoa(l3),
		strconv.Itoa(line),
		"1MB",
		strconv.Itoa(n),
	}
	results = append(results, result)

	// Also compare for slightly smaller than cache size to account for randomness
	runTime = measurePerformance(func(a, b [][]int) [][]int {
		return matmul.CacheAwareMM(a, b, matmul.ComputeBlockSize("900KB"))
	}, pairs)
	log.Println("Cache aware matrix multiplication speedup (900KB): ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))
	result = []string{
		time.Now().Format("2006-01-02 15:04:05"),
		"MM-CacheAware",
		strconv.Itoa(dim1),
		strconv.Itoa(dim2),
		strconv.FormatInt(runTime.Nanoseconds(), 10),
		strconv.FormatFloat(float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()), 'f', -1, 64),
		strconv.Itoa(numProc),
		strconv.Itoa(l1),
		strconv.Itoa(l2),
		strconv.Itoa(l3),
		strconv.Itoa(line),
		"900KB",
		strconv.Itoa(n),
	}
	results = append(results, result)

	// All checks passed
	log.Println("All Matrix multiplication checks passed.")

	// -----------------------------------------------------
	// Merge-Sort Algorithms
	// -----------------------------------------------------

	// Sequential Merge Sort
	d := matrix.MagicMatrix1D(1000000, 42)
	e := mergesort.MergeSort(d)
	if !mergesort.CheckSorted(d, e) {
		panic("Seq. Merge Sort check failed")
	}
	baseline = measurePerformance1D(mergesort.MergeSort, arrays)
	log.Println("Seq. Merge Sort running time: ", baseline.Nanoseconds())
	log.Println("Seq. Merge Sort check passed.")
	result = []string{
		time.Now().Format("2006-01-02 15:04:05"),
		"MS-Seq",
		strconv.Itoa(dim1),
		"",
		strconv.FormatInt(baseline.Nanoseconds(), 10),
		strconv.Itoa(1),
		strconv.Itoa(numProc),
		strconv.Itoa(l1),
		strconv.Itoa(l2),
		strconv.Itoa(l3),
		strconv.Itoa(line),
		"",
		strconv.Itoa(n),
	}
	results = append(results, result)

	// Parallel Merge sort
	e = mergesort.ParMergeSort(d)
	if !mergesort.CheckSorted(d, e) {
		panic("Par Merge Sort check failed")
	}
	runTime = measurePerformance1D(mergesort.ParMergeSort, arrays)
	log.Println("Par Merge Sort speedup: ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))
	log.Println("Par Merge Sort check passed.")
	result = []string{
		time.Now().Format("2006-01-02 15:04:05"),
		"MS-Par",
		strconv.Itoa(dim1),
		"",
		strconv.FormatInt(runTime.Nanoseconds(), 10),
		strconv.FormatFloat(float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()), 'f', -1, 64),
		strconv.Itoa(numProc),
		strconv.Itoa(l1),
		strconv.Itoa(l2),
		strconv.Itoa(l3),
		strconv.Itoa(line),
		"",
		strconv.Itoa(n),
	}
	results = append(results, result)

	// Parallel Merge Sort using PMerge (Binary Search)
	e = mergesort.PmergeSort(d)
	if !mergesort.CheckSorted(d, e) {
		panic("P Merge Sort check failed")
	}
	runTime = measurePerformance1D(mergesort.PmergeSort, arrays)
	log.Println("P Merge Sort speedup: ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))
	log.Println("P Merge Sort check passed.")
	result = []string{
		time.Now().Format("2006-01-02 15:04:05"),
		"MS-PMerge",
		strconv.Itoa(dim1),
		"",
		strconv.FormatInt(runTime.Nanoseconds(), 10),
		strconv.FormatFloat(float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()), 'f', -1, 64),
		strconv.Itoa(numProc),
		strconv.Itoa(l1),
		strconv.Itoa(l2),
		strconv.Itoa(l3),
		strconv.Itoa(line),
		"",
		strconv.Itoa(n),
	}
	results = append(results, result)

	// -----------------------------------------------------
	// Store Results
	// -----------------------------------------------------
	err := storeResults("./data/"+filename+".csv", results)
	if err != nil {
		log.Println("❌ Failed to write:", err)
	} else {
		log.Println("✅ Appended results to CSV.")
	}

	// ------------------------------------------------------
	// Funnelsort
	// ------------------------------------------------------

	// TestSpeedup()
	// TestCorrectness()
}
