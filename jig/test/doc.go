// Package test provides examples, tests and benchmarks for the channel (specialized on type int).
//
// To run benchmarks for the channel several times, use the following command:
// 	jig -r && go test -run=XXX -bench=Chan -cpu=1,2,3,4,5,6,7,8 -timeout=1h -count=10
package test

import _ "github.com/reactivego/channel/jig"

const BUFSIZE = 512
