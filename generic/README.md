# multicast

    import "github.com/reactivego/multicast/generic"

[![](https://godoc.org/github.com/reactivego/multicast/generic?status.png)](http://godoc.org/github.com/reactivego/multicast/generic)

Package `multicast` implements a generic multicast channel. This channel supports multiple concurrent senders and multicasting to more than one receiver. Additionally, code that receives from the channel can specify to receive only messages younger than a given maximum age. Install [Generics for Go](https://github.com/reactivego/jig/) to use the library.
