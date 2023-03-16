package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"time"
)

type Items struct {
	items []struct {
		Name string
		Age  int
	}
}

func main() {
	sc, _ := stan.Connect("test-cluster", "subscriber")

	// Simple Async Subscriber
	sub, err := sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	fmt.Println(err)

	// Simple Synchronous Publisher
	err = sc.Publish("foo", []byte("Hello World")) // does not return until an ack has been received from NATS Streaming
	fmt.Println(err)
	time.Sleep(time.Second)

	// Unsubscribe
	sub.Unsubscribe()

	// Close connection
	sc.Close()
}
