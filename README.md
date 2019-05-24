# multicast

    import "github.com/reactivego/multicast"

[![](https://godoc.org/github.com/reactivego/multicast?status.png)](http://godoc.org/github.com/reactivego/multicast)

Package `multicast` offers a `Chan` type that can multicast and replay the messages you send to it to multiple receivers.

The standard Go channel cannot multicast the same message to multiple receivers and it cannot play back messages previously sent to it. The `Chan` type offered here allows multicasting and playback of messages. You can even limit playback to messages younger than a certain age because `Chan` also stores a timestamp with each message send.

This multicast channel is different from other multicast implementations in that it uses only fast synchronization primitives like atomic operations to implement its features. Furthermore, it also doesn't use goroutines internally. This allows it to operate at a very high level of performance.

If you are in a situation where you need to record and replay a stream of data or need to split a stream of data into multiple identical streams, then this package offers a fast and simple implementation.

## License
This library is licensed under the terms of the MIT License. See [LICENSE](LICENSE) file in this repository for copyright notice and exact wording.
