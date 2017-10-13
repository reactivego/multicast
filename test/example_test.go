package test

import "fmt"

func Example_simple() {

	ch := NewChanInt(32, 8)

	// Send
	for i := 0; i < 5; i++ {
		ch.Send(i)
	}
	ch.Close(nil)

	// Receive
	ep, _ := ch.NewEndpoint(ReplayAll)
	ep.Range(func(value int, err error, closed bool) bool {
		if !closed {
			fmt.Printf("%d ", value)
		} else {
			fmt.Print("closed")
		}
		return true
	}, 0)

	// Output:
	// 0 1 2 3 4 closed
}
