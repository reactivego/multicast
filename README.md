# channel

    import "github.com/reactivego/channel"

[![](https://godoc.org/github.com/reactivego/channel?status.png)](http://godoc.org/github.com/reactivego/channel)

Package `channel` provides a multi-sender merging&multicasting channel. It is a specialization of the generics library
github.com/reactivego/channel/jig on the emtpy interface type.

The channel can be used by multiple senders to simultaneously send messages to the channel that get merged and buffered and are then multicasted to multiple concurrent receivers. There is support for a replay buffer to replay messages received in the past to newly connecting receivers. Messages are timestamped, so during receiving the maximum age of the messages will cause old messages to be skipped.

## License
This library is licensed under the terms of the MIT License. See [LICENSE](LICENSE) file in this repository for copyright notice and exact wording.
