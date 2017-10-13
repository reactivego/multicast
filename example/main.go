package main

import (
	"fmt"
	"time"

	_ "github.com/reactivego/channel"
)

//jig:file {{.package}}.go

func main() {
	for i := 0; i < 10; i++ {
		FastChanTimedBurst1S8R()
		GoChanTimedBurst1S8R()
	}
}

const BUFSIZE = 1024
const MAXPAR = 8

func FastChanTimedBurst1S8R() {
	drain := func(e *EndpointInt, count *int32, wait chan struct{}) {
		e.Range(func(next int, err error, done bool) bool {
			(*count)++
			return true
		})
		close(wait)
	}

	c := NewChanInt(BUFSIZE, MAXPAR)
	if c == nil {
		println("can't create chanel")
		return
	}

	var counters []*int32
	var endpoints []*EndpointInt
	var waits []chan struct{}
	for i := 0; i < MAXPAR; i++ {
		ep, err := c.NewEndpoint(ReplayAll)
		if err != nil {
			println(err)
			return
		}
		count := int32(0)
		wait := make(chan struct{})
		go drain(ep, &count, wait)
		endpoints = append(endpoints, ep)
		counters = append(counters, &count)
		waits = append(waits, wait)
	}

	count := int32(0)
	end := time.Now().Add(time.Second)
	for i := 0; time.Now().Before(end); i++ {
		c.FastSend(i)
		count++
	}
	c.Close(nil)

	for p := 0; p < MAXPAR; p++ {
		<-waits[p]
	}

	tps := float64(1000000000) / float64(count)
	fmt.Printf("<fc: %.0f ns/send %.1fM msgs/sec\n", tps, float64(count)/1000000.0)
	for p := 0; p < MAXPAR; p++ {
		tps := float64(1000000000) / float64(*counters[p])
		fmt.Printf("%d: %.0f ns/send %.1fM msgs/sec\n", p, tps, float64(*counters[p])/1000000.0)
	}
}

func GoChanTimedBurst1S8R() {
	drain := func(ch chan interface{}, count *int32, wait chan struct{}) {
		for range ch {
			(*count)++
		}
		close(wait)
	}

	var counters []*int32
	var channels []chan interface{}
	var waits []chan struct{}

	for i := 0; i < MAXPAR; i++ {
		ch := make(chan interface{}, BUFSIZE)
		count := int32(0)
		wait := make(chan struct{})
		go drain(ch, &count, wait)
		channels = append(channels, ch)
		counters = append(counters, &count)
		waits = append(waits, wait)
	}

	count := int32(0)
	end := time.Now().Add(time.Second)
	for i := 0; time.Now().Before(end); i++ {
		for p := 0; p < MAXPAR; p++ {
			channels[p] <- i
		}
		count++
	}
	for p := 0; p < MAXPAR; p++ {
		close(channels[p])
	}

	for p := 0; p < MAXPAR; p++ {
		<-waits[p]
	}

	tps := float64(1000000000) / float64(count)
	fmt.Printf("<go: %.0f ns/send %.1fM msgs/sec\n", tps, float64(count)/1000000.0)
	for p := 0; p < MAXPAR; p++ {
		tps := float64(1000000000) / float64(*counters[p])
		fmt.Printf("%d: %.0f ns/send %.1fM msgs/sec\n", p, tps, float64(*counters[p])/1000000.0)
	}
}
