package main

import (
	"sync"
)

const N = 1024

var (
	A   = [N][N]int{}
	B   = [N][N]int{}
	Res = [N][N]int{}
)

func doSomeCPUIntensiveStuff() {
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			A[i][j] = 10
			B[i][j] = 20
			Res[i][j] = 0
		}
	}

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			for k := 0; k < N; k++ {
				Res[i][k] += A[i][j] * B[j][k]
			}
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			doSomeCPUIntensiveStuff()
		}()
	}
	wg.Wait()
}
