package main

import (
	"fmt"
	"sync"
)

// producer - will produce the jobs
// consumer - will consume the jobs
// channel - acts as a queue for this transportation

// ch channel sends int values out ch chan <- int means it is send-only channel
// this function can write values to the channel but cannot read the values
func producer(id int, ch chan<- int) {
	for i := 0; i <= 100; i++ {
		// sends the value i in the channel and since this is unbuff channel it blocks here until someone recieve on it
		ch <- i
		fmt.Printf("producer %d sent %d\n", id, i)
	}
}

// consumer continously receives values from the channels
// <-chan means receive-only channel
func consumer(id int, ch <-chan int) {
	// range keeps receiving values until
	// 1. channel is closed or 2. all buffered values are drained
	// once the channel is closed or empty the loop exits
	for val := range ch {
		fmt.Printf("consumer %d received %d\n", id, val)
	}
	// only reaches here when the channel closed
	fmt.Printf("consumer %d exited", id)
}

func main() {
	// shared communication channel
	// producers send value into channel and consumers consume values from it
	// creating unbuffered channel
	ch := make(chan int)
	var pwg sync.WaitGroup
	var cwg sync.WaitGroup

	numProducers := 2
	numConsumers := 3

	// start the producers
	for i := 1; i <= numProducers; i++ {
		pwg.Add(1)

		go func(id int) {
			// must always call wg.Done() when goroutines finishes
			defer pwg.Done()
			producer(id, ch)
		}(i)
	}

	// start the consumers 
	for i := 1; i <= numConsumers; i++ {
		cwg.Add(1)

		go func(id int) {
			defer cwg.Done()
			consumer(id, ch)
		}(i)
	}

	// channel closer goroutine
	// hard and fast rule: only sender side can close the channel because
	// Consumers dont know how many producers exist and how many more values are coming
	// therefore this goroutine waits until all producers are done and then closes the channel exactly once
	// 
	go func() {
		pwg.Wait()
		fmt.Println("all producers are done")
		// now no more values will be sent into channel
		close(ch)
	}()

	// consumer continues proccesing until channel closes and all remaining values are drained
	cwg.Wait()
	fmt.Println("all consumers finished")
	
}