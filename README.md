# channel

    import "github.com/reactivego/channel"

[![](https://godoc.org/github.com/reactivego/channel?status.png)](http://godoc.org/github.com/reactivego/channel)

Library `channel` provides a multisender, multicasting buffered channel for [Go](https://golang.org/). It is a generics library for creating type safe channels that are much more versatile than the native channels provided by Go.

Unlike native Go channels, messages send to a `reactivego/channel` are multicasted to all receivers. A new endpoint created while the channel is operational can choose to receive messages previously sent by specifying a replay count parameter, or 0 to indicate it is only interested in new messages.

Just like native Go channels, the channel exhibits blocking backpressure to the sender goroutines when the channel buffer is full. Total speed of the channel is dictated by the slowest receiver.

This is a 'Just-In-time Generics' library for Go. The way in which a channel
is created will determine strong or weak typing. Channels can be strongly
typed by specifying an explicit type, for example:

```go
NewChanInt(128,8)
NewChanString(128,8)
```

Or alternatively you can send heterogeneous messages on an `interface{}` typed
channel created as follows:

```go
NewChan(128,8)
```

Because the generics in this library are only recognized by [Just-in-time Generics for Go](https://github.com/reactivego/jig/), you will need to install it.

## License
This library is licensed under the terms of the MIT License. See [LICENSE](LICENSE) file in this repository for copyright notice and exact wording.
