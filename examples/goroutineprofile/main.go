package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rakyll/autopprof"
)

func main() {

	// taken from the following example: https://talks.golang.org/2012/concurrency.slide#27
	c := fanIn(boring("Joe"), boring("Ann"))

	for i := 0; i < 10; i++ {
		fmt.Println(<-c) // display any message received on the FanIn channel
	}

	fmt.Println("You're boring: I'm leaving")

	autopprof.Capture(autopprof.GoRoutineProfile{})

	time.Sleep(5 * time.Second)

}

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string) // The FanIn channel

	go func() {
		for {
			c <- <-input1 // write the message to the FanIn channel, Blocking call.
		}
	}()

	go func() {
		for {
			c <- <-input2 // write the message to the FanIn channel, Blocking call
		}
	}()

	return c
}

// returns receive-only (<-) channel of string.
func boring(msg string) <-chan string {
	c := make(chan string)

	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()

	return c
}
