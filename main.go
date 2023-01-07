package main

import "fmt"

func main() {
	msg := make(chan string)
	go func() {
		msg <- "Hello nyse & nasdaq"
	}()
	mx := <-msg
	fmt.Printf("hello stock market!: %v", mx)
}
