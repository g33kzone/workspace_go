package main

import (
	"fmt"
	"time"
)

func main() {

	c1 := make(chan int)
	c2 := make(chan int)

	go func() {
		time.Sleep(time.Second * 1)
		c1 <- add(1, 2)
	}()
	go func() {
		time.Sleep(time.Second * 2)
		c2 <- sub(1, 2)
	}()

	for i := 0; i < 2; i++ {
		// Await both of these values
		// simultaneously, printing each one as it arrives.
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
	fmt.Println("Test Print...")
}

func add(num1 int, num2 int) int {
	return num1 + num2
}

func sub(num1 int, num2 int) int {
	return num2 - num1
}
