// Library channel provides a (jig) multicasting, multi sender/receiver,
// buffered channel.
//
// This is a 'Just-In-time Generics' library for Go. Channels can be strongly
// typed using e.g. NewChanInt() or NewChanString(). Or alternatively you can
// send heterogeneous messages on an interface{} typed channel using NewChan().
//
// Because the generics definitions in 'channel' are only recognized by
// Just-in-time Generics for Go, you will need to install the jig tool
// (https://github.com/reactivego/jig/).
//
// Unlike native Go channels, messages send to this channel are multicasted to
// all receivers. A new endpoint created while the channel is operational can
// choose to receive messages previously sent by specifying a keep parameter,
// or 0 to indicate it is only interested in new messages.
//
// Just like native Go channels, the channel exhibits blocking backpressure to
// the sender goroutines when the channel buffer is full. Total speed of the
// channel is dictated by the slowest receiver.
//
// Benchmarks
//
// Note that benchmarks should always be taken with a grain of salt, but since
// I'm not very well versed in doing Go benchmarks, I might have messed
// something up. So if you see obvious flaws in my method of testing, please let
// me know so I can improve it. As it is, I find the results I found for my
// implementation freakishly fast. Started out 4 times slower than Go though so
// it's not as if I did not have to work for it.
//
// Initially I used interface{} as the message type. Then, later switched to
// int after I converted the library to a generics library and I could generate
// for any type. Speedwise there is not much difference between using either
// int or interface{}. Benchmark results are for the int type though.
//
// Fan-out
//
// A fan-out configuration (buffer capacity 512) where a single sender is
// transmitting int values to multiple receivers where messages are multicasted
// so every receiver receives the same set of messages performs as follows:
//
// 	Mac:test $ go test -run=XXX -bench=FastChan.FanOut -cpu=1,2,3,4,5,6,7,8
// 	goos: darwin
// 	goarch: amd64
// 	pkg: github.com/reactivego/channel/test
// 	BenchmarkFastChan_FanOut_1S_NR     	50000000	        27.6 ns/op
// 	BenchmarkFastChan_FanOut_1S_NR-2   	100000000	        13.9 ns/op
// 	BenchmarkFastChan_FanOut_1S_NR-3   	200000000	         7.96 ns/op
// 	BenchmarkFastChan_FanOut_1S_NR-4   	100000000	        20.6 ns/op
// 	BenchmarkFastChan_FanOut_1S_NR-5   	200000000	        12.9 ns/op
// 	BenchmarkFastChan_FanOut_1S_NR-6   	200000000	        12.8 ns/op
// 	BenchmarkFastChan_FanOut_1S_NR-7   	100000000	        10.7 ns/op
// 	BenchmarkFastChan_FanOut_1S_NR-8   	100000000	        10.4 ns/op
// 	PASS
// 	ok  	github.com/reactivego/channel/test	15.782s
//
// From the results we can see that even for 8 concurrent receivers, the
// receiving goroutines are not contending for access to the data.
//
// The same configuration, but implemented using multiple parallel native Go
// channels. Since Go doesn't support multicasting to multiple receivers from
// a single channel. Our multiple channels assembly gives the following result:
//
// 	Mac:test $ go test -run=XXX -bench=GoChan.FanOut -cpu=1,2,3,4,5,6,7,8
// 	goos: darwin
// 	goarch: amd64
// 	pkg: github.com/reactivego/channel/test
// 	BenchmarkGoChan_FanOut_1S_NR     	20000000	       100 ns/op
// 	BenchmarkGoChan_FanOut_1S_NR-2   	 5000000	       263 ns/op
// 	BenchmarkGoChan_FanOut_1S_NR-3   	 3000000	       542 ns/op
// 	BenchmarkGoChan_FanOut_1S_NR-4   	 1000000	      1499 ns/op
// 	BenchmarkGoChan_FanOut_1S_NR-5   	 1000000	      1739 ns/op
// 	BenchmarkGoChan_FanOut_1S_NR-6   	 1000000	      2146 ns/op
// 	BenchmarkGoChan_FanOut_1S_NR-7   	  500000	      2498 ns/op
// 	BenchmarkGoChan_FanOut_1S_NR-8   	  500000	      2909 ns/op
// 	PASS
// 	ok  	github.com/reactivego/channel/test	14.093s
//
// What we see here is the sender slowing down as it is pumping the same
// information into increasingly more separate channels.
//
// Fan-in
//
// A fan-in configuration (buffer capacity 512) where multiple senders are
// transmitting int values to a single receiver. Messages are merged in
// arbitrary order and all delivered to the receiver.
//
// 	Mac:test $ go test -run=XXX -bench=FastChan.FanIn -cpu=1,2,3,4,5,6,7,8
// 	goos: darwin
// 	goarch: amd64
// 	pkg: github.com/reactivego/channel/test
// 	BenchmarkFastChan_FanIn_NS_1R     	50000000	        33.6 ns/op
// 	BenchmarkFastChan_FanIn_NS_1R-2   	50000000	        29.9 ns/op
// 	BenchmarkFastChan_FanIn_NS_1R-3   	30000000	        51.2 ns/op
// 	BenchmarkFastChan_FanIn_NS_1R-4   	30000000	        46.8 ns/op
// 	BenchmarkFastChan_FanIn_NS_1R-5   	30000000	        44.5 ns/op
// 	BenchmarkFastChan_FanIn_NS_1R-6   	30000000	        44.1 ns/op
// 	BenchmarkFastChan_FanIn_NS_1R-7   	30000000	        44.1 ns/op
// 	BenchmarkFastChan_FanIn_NS_1R-8   	30000000	        45.8 ns/op
// 	PASS
// 	ok  	github.com/reactivego/channel/test	11.832s
//
// I really had to work hard to get performance to an acceptable level. Started
// out an order of magnitude slower than native Go, as the ammount of contention
// for goroutines trying to gain write access to the channel was insane.
// Eventually I just changed the solution to hand out write slots to the
// concurrent senders and have one of the receiver goroutines consolidate and
// commit the data by looking for contiguous sequences of slots marked as
// updated by their sender goroutines. So data is ordered on slot hand-out time
// not on actual data write time. But, that fits with the semantics you'd expect
// of a concurrent sender channel, so all in all it's a good approach.
//
// 	Mac:test $ go test -run=XXX -bench=GoChan.FanIn -cpu=1,2,3,4,5,6,7,8
// 	goos: darwin
// 	goarch: amd64
// 	pkg: github.com/reactivego/channel/test
// 	BenchmarkGoChan_FanIn_NS_1R     	20000000	        60.3 ns/op
// 	BenchmarkGoChan_FanIn_NS_1R-2   	20000000	        81.2 ns/op
// 	BenchmarkGoChan_FanIn_NS_1R-3   	20000000	        95.3 ns/op
// 	BenchmarkGoChan_FanIn_NS_1R-4   	20000000	       108 ns/op
// 	BenchmarkGoChan_FanIn_NS_1R-5   	10000000	       131 ns/op
// 	BenchmarkGoChan_FanIn_NS_1R-6   	10000000	       160 ns/op
// 	BenchmarkGoChan_FanIn_NS_1R-7   	10000000	       182 ns/op
// 	BenchmarkGoChan_FanIn_NS_1R-8   	10000000	       199 ns/op
// 	PASS
// 	ok  	github.com/reactivego/channel/test	14.714s
//
// Go natively supports fan-in so its performance was quite good!
package channel

//jig:file support.go

type foo interface{}
