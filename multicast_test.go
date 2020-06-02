package multicast_test

import (
	"fmt"
	"sync"

	"github.com/reactivego/multicast"
)

func Example_fastSend1x2() {
	ch := multicast.NewChan(128, 2)

	// FastSend allows only a single goroutine sending and does not store
	// timestamps with messages.

	ch.FastSend("Hello")
	ch.FastSend("World!")
	ch.Close(nil)
	if ch.Closed() {
		fmt.Println("channel closed")
	}

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

	var wg sync.WaitGroup
	wg.Add(2)
	ep1, _ := ch.NewEndpoint(multicast.ReplayAll)
	go func() {
		ep1.Range(print, 0)
		wg.Done()
	}()

	ep2, _ := ch.NewEndpoint(multicast.ReplayAll)
	go func() {
		ep2.Range(print, 0)
		wg.Done()
	}()
	wg.Wait()

	// Unordered Output:
	// channel closed
	// Hello
	// Hello
	// World!
	// World!
	// closed
	// closed
}

func Example_send2x2() {
	ch := multicast.NewChan(128, 2)

	// Send suppports multiple goroutine sending and stores a timestamp with
	// every message sent.

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		ch.Send("Hello")
		wg.Done()
	}()
	go func() {
		ch.Send("World!")
		wg.Done()
	}()
	wg.Wait()
	ch.Close(nil)

	if ch.Closed() {
		fmt.Println("channel closed")
	}

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

	wg.Add(2)
	ep1, _ := ch.NewEndpoint(multicast.ReplayAll)
	go func() {
		ep1.Range(print, 0)
		wg.Done()
	}()

	ep2, _ := ch.NewEndpoint(multicast.ReplayAll)
	go func() {
		ep2.Range(print, 0)
		wg.Done()
	}()
	wg.Wait()

	// Unordered Output:
	// channel closed
	// Hello
	// Hello
	// World!
	// World!
	// closed
	// closed
}
