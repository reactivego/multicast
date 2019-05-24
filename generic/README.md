# multicast

    import "github.com/reactivego/multicast/generic"

[![](https://godoc.org/github.com/reactivego/multicast/generic?status.png)](http://godoc.org/github.com/reactivego/multicast/generic)

Package `generic` is a Just-in-time Generics (jig) library. It contains a generic implementation of a multicast channel. This supports multiple concurrent senders and also multicasting to more than one receiver. Additionally code that receives from the channel can specify to receive only messages younger a given maximum age. Install [Just-in-time Generics for Go](https://github.com/reactivego/jig/) to use the library.
