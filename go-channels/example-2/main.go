package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	Find()
	// time.Sleep(10 * time.Second)
}

// Response -
type Response struct {
	data   interface{}
	status bool
}

func cntxt() context.Context {
	return context.Background()
}

// Find -
func Find() {
	// ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(cntxt(), time.Second*2)
	ch := make(chan Response, 1)

	go func() {
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			fmt.Println("Canceled by timeout")
			return
		}

		// ch <- Response{data: "data", status: true}
	}()

	select {
	case <-ch:
		fmt.Println("Read from ch")
	case <-time.After(time.Second * 5):
		fmt.Println("Timed out")
		cancel()
	}

}
