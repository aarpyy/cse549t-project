// run with command: go test -bench=.


package main

import (
    "math"
    "math/rand"
    "testing"
)

// cache aware
func blockedMatrixMultiply(A, B, C [][]int, n, blockSize int) {
    for i := 0; i < n/blockSize; i++ {
        for j := 0; j < n/blockSize; j++ {
            for k := 0; k < n/blockSize; k++ {
                multiplySubMatrix(A, B, C, i, j, k, blockSize)
            }
        }
    }
}


func multiplySubMatrix(A, B, C [][]int, i, j, k, blockSize int) {
    rowStart := i * blockSize
    colStart := j * blockSize
    sharedStart := k * blockSize

    for ii := rowStart; ii < rowStart+blockSize; ii++ {
        for jj := colStart; jj < colStart+blockSize; jj++ {
            for kk := sharedStart; kk < sharedStart+blockSize; kk++ {
                C[ii][jj] += A[ii][kk] * B[kk][jj]
            }
        }
    }
}

func randMatrix(n int) [][]int {
    m := make([][]int, n)
    for i := range m {
        m[i] = make([]int, n)
        for j := range m[i] {
            m[i][j] = rand.Intn(10)
        }
    }
    return m
}



// cahce oblivious
func cacheObliviousMatrixMultiply(A, B, C [][]int, aRow, aCol, bRow, bCol, cRow, cCol, size int) {
    if size == 1 {
        C[cRow][cCol] += A[aRow][aCol] * B[bRow][bCol]
        return
    }

    half := size / 2

    cacheObliviousMatrixMultiply(A, B, C, aRow, aCol, bRow, bCol, cRow, cCol, half)
    cacheObliviousMatrixMultiply(A, B, C, aRow, aCol+half, bRow+half, bCol, cRow, cCol, half)

    cacheObliviousMatrixMultiply(A, B, C, aRow, aCol, bRow, bCol+half, cRow, cCol+half, half)
    cacheObliviousMatrixMultiply(A, B, C, aRow, aCol+half, bRow+half, bCol+half, cRow, cCol+half, half)

    cacheObliviousMatrixMultiply(A, B, C, aRow+half, aCol, bRow, bCol, cRow+half, cCol, half)
    cacheObliviousMatrixMultiply(A, B, C, aRow+half, aCol+half, bRow+half, bCol, cRow+half, cCol, half)

    cacheObliviousMatrixMultiply(A, B, C, aRow+half, aCol, bRow, bCol+half, cRow+half, cCol+half, half)
    cacheObliviousMatrixMultiply(A, B, C, aRow+half, aCol+half, bRow+half, bCol+half, cRow+half, cCol+half, half)
}


func subMatrix(M [][]int, rowStart, colStart, size int) [][]int {
    sub := make([][]int, size)
    for i := 0; i < size; i++ {
        sub[i] = M[rowStart+i][colStart : colStart+size]
    }
    return sub
}


// testing
func BenchmarkMatrixAlgorithms(b *testing.B) {
    n := 512

    A := randMatrix(n)
    B := randMatrix(n)

    COblivious := make([][]int, n)
    for i := range COblivious {
        COblivious[i] = make([]int, n)
    }

    b.Run("CacheOblivious", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            for x := range COblivious {
                for y := range COblivious[x] {
                    COblivious[x][y] = 0
                }
            }
            cacheObliviousMatrixMultiply(A, B, COblivious, 0, 0, 0, 0, 0, 0, n)
        }
    })

    z := 32768
    elementSize := 8
    zElements := z / elementSize
    blockSize := int(math.Sqrt(float64(zElements)) / 3)

    CAware := make([][]int, n)
    for i := range CAware {
        CAware[i] = make([]int, n)
    }

    b.Run("CacheAware", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            for x := range CAware {
                for y := range CAware[x] {
                    CAware[x][y] = 0
                }
            }
            blockedMatrixMultiply(A, B, CAware, n, blockSize)
        }
    })
}
