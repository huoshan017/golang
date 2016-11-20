package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	defer func() {
		e := recover()
		if e != nil {
			log.Println(e)
		}
		var inp int
		fmt.Scanln(&inp)
	}()

	a := []int{11, 43, 4, 5, 7, 3, 32, 123, 9, 21, 32,
		41, 53, 12, 15, 76, 29, 65, 111, 234, 123,
		321, 124, 22, 33, 44, 88, 55, 77}

	//b := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	start := time.Now()
	fmt.Printf("begin time: %v\n", start.UnixNano())
	Qsort(0, a, 0, len(a)-1)
	end := time.Now()
	fmt.Printf("end time: %v,  cost %v\n", end.UnixNano(), time.Since(start))
	fmt.Printf("arr=%v\n", a)
	start = time.Now()
	fmt.Printf("begin time: %v\n", start.UnixNano())
	Qsort(1, a, 0, len(a)-1)
	end = time.Now()
	fmt.Printf("end time: %v,  cost %v\n", end.UnixNano(), time.Since(start))
	fmt.Printf("arr=%v  end!!!\n", a)
}
