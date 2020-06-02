package multicast_test

import (
	"fmt"

	"github.com/reactivego/multicast"
)

func Example() {
	ch := multicast.NewChan(128, 1)
	if true {
		ch.FastSend("hello")
	} else {
		ch.Send("world!")
	}
	ch.Close(nil)
	if ch.Closed() {
		fmt.Println("channel closed")
	}

	ep, _ := ch.NewEndpoint(multicast.ReplayAll)
	print := func(value interface{}, err error, closed bool) bool {
		switch {
		case !closed:
			fmt.Println(value)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("closed")
		}
		return true
	}
	ep.Range(print, 0)
	ep.Cancel()

	// Output:
	// channel closed
	// hello
	// closed
}
