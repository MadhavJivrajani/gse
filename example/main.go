package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"golang.org/x/sys/unix"
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

func doSomeModeratelyCPUIntensiveStuff() {
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			A[i][j] = 10
			B[i][j] = 20
			Res[i][j] = 0
		}
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			go func() {
				defer wg.Done()
				runtime.LockOSThread()
				// defer runtime.UnlockOSThread() !!!
				tid := unix.Gettid()
				fmt.Println("[Intensive] TID:", tid)
				fmt.Fprintf(os.Stderr, "%d\n", tid)
				doSomeCPUIntensiveStuff()
			}()
		} else {
			go func() {
				defer wg.Done()
				runtime.LockOSThread()
				// defer runtime.UnlockOSThread() !!!
				fmt.Println("[Moderately Intensive] TID:", unix.Gettid())
				doSomeModeratelyCPUIntensiveStuff()
			}()
		}
	}
	wg.Wait()
}

// 				_, err := unix.Getpriority(unix.PRIO_PROCESS, tid)
// 				if err != nil {
// 					log.Fatal(err)
// 				}
// 				err = unix.Setpriority(unix.PRIO_PROCESS, tid, -20)
// 				if err != nil {
// 					log.Fatal(err)

// 				}

// go func(i int) {
// 	defer wg.Done()
// 	runtime.LockOSThread()
// 	// defer runtime.UnlockOSThread() !!!
// 	log.Println("[Intensive] TID:", unix.Gettid())
// 	if i == 4 {
// 		tid := unix.Gettid()
// 		log.Println("[Intensive][High Priority] TID:", tid)
// 		_, err := unix.Getpriority(unix.PRIO_PROCESS, tid)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		err = unix.Setpriority(unix.PRIO_PROCESS, tid, -20)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	doSomeCPUIntensiveStuff()
// }(i)
