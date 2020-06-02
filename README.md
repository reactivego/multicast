# multicast

    import "github.com/reactivego/multicast"

[![](svg/godev.svg)](https://pkg.go.dev/github.com/reactivego/multicast?tab=doc)
[![](svg/godoc.svg)](https://godoc.org/github.com/reactivego/multicast)

Package `multicast` provides MxN multicast channels for Go with buffering and time based buffer eviction.
It can be fed by multiple concurrent senders. It multicasts and replays messages to multiple concurrent receivers.

If you are in a situation where you need to record and replay a stream of data or need to split a stream of data into multiple identical streams, then this package offers a fast and simple implementation.

## Example
Code:
```go
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
```
Output:
```
channel closed
hello
closed
```

## Compared to Go channels
The standard Go channel cannot multicast the same message to multiple receivers and it cannot play back messages previously sent to it. The `multicast.Chan` type offered here does.

Additionally, you can even evict messages from the buffer that are past a certain age because `multicast.Chan` also stores a timestamp with each message sent.

## Compared to other Multicast packages
This multicast channel is different from other multicast implementations.

1. It uses only fast synchronization primitives like atomic operations to implement its features.
2. It doesn't use goroutines internally.
3. It uses internal struct padding to speed up CPU cache access.

This allows it to operate at a very high level of performance.

## Regenerating this Package
This package is generated from generics in the sub-folder `generic` by the [jig](http://github.com/reactivego/jig) tool.
You don't need to regenerate this package in order to use it. However, if you are interested in regenerating it, then read on.

The [jig](http://github.com/reactivego/jig) tool provides the parametric polymorphism capability that Go 1 is missing.
It works by replacing place-holder types of generic functions and datatypes with `interface{}` (it can also generate statically typed code though).

To regenerate, change the current working directory to the package directory and run the [jig](http://github.com/reactivego/jig) tool as follows:

```bash
$ go get -d github.com/reactivego/jig
$ go run github.com/reactivego/jig -v
```

## License
This library is licensed under the terms of the MIT License. See [LICENSE](LICENSE) file in this repository for copyright notice and exact wording.
