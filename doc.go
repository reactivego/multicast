/*
Package multicast provides a Chan type that can multicast and replay messages to
multiple receivers.

Multicast and Replay

Native Go channels don't support multicasting the same message to multiple
receivers and they don't support replaying previously sent messages.

Unlike native Go channels, messages send to this channel are multicast to
all receiving endpoints. A new endpoint created while the channel is operational
can choose to receive messages previously sent by specifying a replay count
parameter, or 0 to indicate it is only interested in new messages.

You can also limit playback to messages younger than a certain age because
the channel stores a timestamp with each message you send to it.

Just like native Go channels, the channel exhibits blocking backpressure to
the sender goroutines when the channel buffer is full. Total speed of the
channel is dictated by the slowest receiver.

Lock and Goroutine free

This multicast channel is different from other multicast implementations in
that it uses only fast synchronization primitives like atomic operations to
implement its features. Furthermore, it also doesn't use goroutines internally.
This implementation is low-latency and has a high throughput.

If you are in a situation where you need to record and replay a stream
of data or you need to split a stream of data into multiple identical streams,
then this package offers a fast and simple solution.

Heterogeneous

Heterogeneous simply means that you can mix types, that is very convenient but
not typesafe. The Chan type provided in this package supports sending and
receiving values of mixed type:

	ch := NewChan(128, 1)
	ch.Send("hello")
	ch.Send(42)
	ch.Send(1.6180)
	ch.Close(nil)

Regenerating this Package

The implementation in this package is generated from a generic implementation
of the Chan type found in the subdirectory "generic" inside this package. By
replacing the place-holder type with "interface{}" a heterogeneous Chan type
is created. To regenerate this channel implementation, run jig inside this
package directory:

	go get -d github.com/reactivego/generics/cmd/jig
	go run github.com/reactivego/generics/cmd/jig -v
*/
package multicast

import _ "github.com/reactivego/multicast/generic"
