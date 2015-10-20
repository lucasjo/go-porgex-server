package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	var intdata int32
	ticker := time.NewTicker(time.Millisecond * 500)
	cnt := atomic.LoadInt32(&intdata)
	go func() {
		for t := range ticker.C {
			fmt.Println("tick at ", t)
			cnt = atomic.AddInt32(&intdata, 1)
			fmt.Printf("load int %v\n", cnt)

		}

	}()

	time.Sleep(time.Millisecond * 1500)
	ticker.Stop()
	fmt.Println("ticker stopped")

}
