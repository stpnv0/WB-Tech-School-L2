package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or2(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}

// реализация через рекурсию
//func or(channels ...<-chan interface{}) <-chan interface{} {
//	switch len(channels) {
//	case 0:
//		return nil
//	case 1:
//		return channels[0]
//	}
//
//	res := make(chan interface{})
//
//	go func() {
//		defer close(res)
//		switch len(channels) {
//		case 2:
//			select {
//			case <-channels[0]:
//			case <-channels[1]:
//			}
//		default:
//			select {
//			case <-channels[0]:
//			case <-channels[1]:
//			case <-or(channels[2:]...):
//			}
//		}
//	}()
//
//	return res
//}

// реализация через sync.Once
func or2(channels ...<-chan interface{}) <-chan interface{} {
	res := make(chan interface{})

	var s sync.Once
	for _, c := range channels {
		go func() {
			<-c

			s.Do(func() {
				close(res)
			})
		}()
	}

	return res
}
