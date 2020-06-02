// Package multicast provides generic MxN multicast channels for Go with
// buffering and time based buffer eviction. It can be fed by multiple
// concurrent senders. It multicasts and replays messages to multiple
// concurrent receivers.
//
// Install the jig tool (https://github.com/reactivego/jig) to use the library.
//
// Unlike native Go channels, messages send to this channel are multicasted to
// all receivers. A new endpoint created while the channel is operational can
// choose to receive messages previously sent by specifying a replay count
// parameter, or 0 to indicate it is only interested in new messages.
//
// Just like native Go channels, the channel exhibits blocking backpressure to
// the sender goroutines when the channel buffer is full. Total speed of the
// channel is dictated by the slowest receiver.
//
// Since this is a generics library, the way in which a channel is created will
// determine strong or weak typing. Channels can be strongly typed by specifying
// an explicit type, for example:
//
//	NewChanInt(128,8)
//	NewChanString(128,8)
//
// Or alternatively you can send heterogeneous messages on an interface{} typed
// channel created as follows:
//
//	 NewChan(128,8)
package multicast

type foo interface{}
