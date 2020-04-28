# test

    import "github.com/reactivego/multicast/test"

[![](https://godoc.org/github.com/reactivego/multicast/test?status.png)](http://godoc.org/github.com/reactivego/multicast/test)

Package `test` provides examples, tests and benchmarks for the multicast channel (specialized on type int).

To run benchmarks for the channel several times, use the following command:

```bash
go run github.com/reactivego/generics/cmd/jig
go test -run=XXX -bench=Chan -cpu=1,2,3,4,5,6,7,8 -timeout=1h -count=10
```
