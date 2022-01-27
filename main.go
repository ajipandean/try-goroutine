package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	o := make(chan int)
	go func() {
		for _, n := range nums {
			o <- n
		}
		close(o)
	}()
	return o
}

func sq(nums <-chan int) <-chan int {
	o := make(chan int)
	go func() {
		for n := range nums {
			o <- n * n
		}
		close(o)
	}()
	return o
}

func mrg(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	o := make(chan int)
	r := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			o <- n
		}
	}
	wg.Add(len(cs))
	for _, n := range cs {
		go r(n)
	}
	go func() {
		wg.Wait()
		close(o)
	}()
	return o
}

func main() {
	in := gen(4, 5, 6, 7, 8, 9)
	r1 := sq(in)
	r2 := sq(in)
	r3 := sq(in)
	for n := range mrg(r1, r2, r3) {
		fmt.Println(n)
	}
}
