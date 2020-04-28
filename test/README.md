# test

    import "github.com/reactivego/multicast/test"

[![](https://godoc.org/github.com/reactivego/multicast/test?status.png)](http://godoc.org/github.com/reactivego/multicast/test)

Package `test` provides examples, tests and benchmarks for the multicast channel (specialized on type int).

To run benchmarks for the channel several times, use the following command:

```bash
go run github.com/reactivego/generics/cmd/jig
go test -run=XXX -bench=Chan -cpu=1,2,3,4,5,6,7,8 -timeout=1h -count=10
```

## Benchmarks

As always, benchmarks should be taken with a grain of salt. When comparing
the performance of our channel implementation against Go's native channels,
there are issues with feature mismatch. These channels have very
different semantics and strenghts. I've tried to create benchmarks that
perform the same amount of work in both implementations.

Initially I used interface{} as the message type. Then, later switched to
int after I converted the library to a generics library and I could generate
for any type. Speedwise there is not much difference between using either
int or interface{}. Benchmark results are for the int type though.

## Fan-out

A fan-out configuration (buffer capacity 512) where a single sender is
transmitting int values to multiple receivers where messages are multicasted
so every receiver receives the same set of messages performs as follows:

```bash
$ go test -run=XXX -bench=FanOut.Chan -cpu=1,2,3,4,5,6,7,8
goos: darwin
goarch: amd64
pkg: github.com/reactivego/multicast/test
BenchmarkFanOut_Chan_1xN     	30000000	        38.3 ns/op
BenchmarkFanOut_Chan_1xN-2   	50000000	        34.7 ns/op
BenchmarkFanOut_Chan_1xN-3   	50000000	        27.3 ns/op
BenchmarkFanOut_Chan_1xN-4   	50000000	        31.6 ns/op
BenchmarkFanOut_Chan_1xN-5   	50000000	        30.4 ns/op
BenchmarkFanOut_Chan_1xN-6   	50000000	        29.5 ns/op
BenchmarkFanOut_Chan_1xN-7   	50000000	        27.9 ns/op
BenchmarkFanOut_Chan_1xN-8   	50000000	        27.2 ns/op
PASS
ok  	github.com/reactivego/multicast/test	11.880s
```

From the results we can see that even for 8 concurrent receivers, the
receiving goroutines are not contending for access to the data.

The same configuration, but implemented using multiple parallel native Go
channels. Since Go doesn't support multicasting to multiple receivers from
a single channel. Our multiple channels assembly gives the following result:

```bash
$ go test -run=XXX -bench=FanOut.Go -cpu=1,2,3,4,5,6,7,8
goos: darwin
goarch: amd64
pkg: github.com/reactivego/multicast/test
BenchmarkFanOut_Go_1xN     	20000000	        61.0 ns/op
BenchmarkFanOut_Go_1xN-2   	20000000	        86.9 ns/op
BenchmarkFanOut_Go_1xN-3   	20000000	        99.8 ns/op
BenchmarkFanOut_Go_1xN-4   	20000000	       115 ns/op
BenchmarkFanOut_Go_1xN-5   	10000000	       182 ns/op
BenchmarkFanOut_Go_1xN-6   	10000000	       197 ns/op
BenchmarkFanOut_Go_1xN-7   	10000000	       204 ns/op
BenchmarkFanOut_Go_1xN-8   	10000000	       210 ns/op
PASS
ok  	github.com/reactivego/multicast/test	16.440s
```

What we see here is the sender slowing down as it is pumping the same
information into increasingly more separate channels.

## Fan-in

A fan-in configuration (buffer capacity 512) where multiple senders are
transmitting int values to a single receiver. Messages are merged in
arbitrary order and all delivered to the receiver.

```bash
$ go test -run=XXX -bench=FanIn.Chan -cpu=1,2,3,4,5,6,7,8
goos: darwin
goarch: amd64
pkg: github.com/reactivego/multicast/test
BenchmarkFanIn_Chan_Nx1     	20000000	        78.0 ns/op
BenchmarkFanIn_Chan_Nx1-2   	20000000	        89.6 ns/op
BenchmarkFanIn_Chan_Nx1-3   	20000000	        75.2 ns/op
BenchmarkFanIn_Chan_Nx1-4   	20000000	        75.5 ns/op
BenchmarkFanIn_Chan_Nx1-5   	20000000	        68.1 ns/op
BenchmarkFanIn_Chan_Nx1-6   	20000000	        65.0 ns/op
BenchmarkFanIn_Chan_Nx1-7   	20000000	        63.7 ns/op
BenchmarkFanIn_Chan_Nx1-8   	20000000	        68.7 ns/op
PASS
ok  	github.com/reactivego/multicast/test	31.698s
```

I really had to work hard to get performance to an acceptable level. Started
out an order of magnitude slower than native Go, as the amount of contention
for goroutines trying to gain write access to the channel was crippling
performance. Eventually, I changed the solution to hand out write slots to the
concurrent senders and have one of the receiver goroutines consolidate and
commit the data by looking for contiguous sequences of slots marked as
updated by their sender goroutines. So data is ordered on slot hand-out time
not on actual data write time. But, that fits with the semantics you'd expect
of a concurrent sender channel, so all in all it's a good approach.

For native Go, the implementation was straight forward as merging message
streams of multiple concurrent senders is standard Go channel functionality.
The results for Go for 1 to 8 concurrent senders and a single receiver are
as follows:

```bash
$ go test -run=XXX -bench=FanIn.Go -cpu=1,2,3,4,5,6,7,8
goos: darwin
goarch: amd64
pkg: github.com/reactivego/multicast/test
BenchmarkFanIn_Go_Nx1     	20000000	        72.9 ns/op
BenchmarkFanIn_Go_Nx1-2   	20000000	       115 ns/op
BenchmarkFanIn_Go_Nx1-3   	20000000	       117 ns/op
BenchmarkFanIn_Go_Nx1-4   	10000000	       133 ns/op
BenchmarkFanIn_Go_Nx1-5   	10000000	       146 ns/op
BenchmarkFanIn_Go_Nx1-6   	10000000	       169 ns/op
BenchmarkFanIn_Go_Nx1-7   	10000000	       184 ns/op
BenchmarkFanIn_Go_Nx1-8   	10000000	       203 ns/op
PASS
ok  	github.com/reactivego/multicast/test	28.924s
```

Go natively supports fan-in, so its performance was very good! However, for
the higher sender counts the performance drops off quite sharply, whereas our
implementation using the 'write slot handout' approach performs much better.

## Fan-In-Out

This benchmark is only implemented for our channel implemenation. It is not
possible to implement this using Go native channels in a very effective way.

What we are benchmarking here is multiple (N) senders concurrently sending
on the channel. The streams of messages are merged into a single stream which
is then multicasted to N concurrent receivers.

```bash
$ go test -run=XXX -bench=FanInOut -cpu=1,2,3,4,5,6,7,8
goos: darwin
goarch: amd64
pkg: github.com/reactivego/multicast/test
BenchmarkFanInOut_Chan_NxN     	20000000	        77.4 ns/op
BenchmarkFanInOut_Chan_NxN-2   	20000000	        99.2 ns/op
BenchmarkFanInOut_Chan_NxN-3   	20000000	       103 ns/op
BenchmarkFanInOut_Chan_NxN-4   	20000000	       101 ns/op
BenchmarkFanInOut_Chan_NxN-5   	20000000	        96.2 ns/op
BenchmarkFanInOut_Chan_NxN-6   	20000000	        93.8 ns/op
BenchmarkFanInOut_Chan_NxN-7   	20000000	        94.8 ns/op
BenchmarkFanInOut_Chan_NxN-8   	20000000	        97.8 ns/op
PASS
ok  	github.com/reactivego/multicast/test	35.091s
```
