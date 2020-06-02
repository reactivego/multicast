// go test -run=XXX -bench=Chan -cpu=1,2,3,4,5,6,7,8 -timeout=1h -count=10

package test

import (
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkFanInOut_Chan_NxN(b *testing.B) {
	PAR := runtime.GOMAXPROCS(0)
	NUM := b.N

	expect := rand.Perm(NUM)
	var count = uint64(NUM)
	var sum = int64(NUM * (NUM - 1) / 2)

	b.ResetTimer()
	start := time.Now()
	channel := NewChanInt(BUFSIZE, PAR)

	var rwgbegin, rwgend sync.WaitGroup
	receiver := func() {
		ep, err := channel.NewEndpoint(ReplayAll)
		rwgend.Add(1)
		rwgbegin.Done()
		if err != nil {
			b.Error(err)
			rwgend.Done()
			return
		}
		var rcount uint64
		var rsum int64
		ep.Range(func(value int, err error, closed bool) bool {
			if !closed {
				rsum += int64(value)
				rcount++
			}
			return true
		}, 0)
		if count != rcount {
			b.Errorf("receiver: rcount(%d) != count(%d)", rcount, count /*, "ep=%+v" ep*/)
		} else if sum != rsum {
			b.Errorf("receiver: rsum(%d) != sum(%d)", rsum, sum)
		}
		rwgend.Done()
	}
	rwgbegin.Add(PAR)
	for i := 0; i < PAR; i++ {
		go receiver()
	}
	rwgbegin.Wait()

	var scount uint64
	var ssum int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := atomic.AddUint64(&scount, 1) - 1
			channel.Send(expect[index])
			atomic.AddInt64(&ssum, int64(expect[index]))
		}
	})
	if sum != ssum {
		b.Errorf("sender: ssum(%d) != sum(%d)", ssum, sum)
	}
	channel.Close(nil)
	rwgend.Wait()

	nps := time.Now().Sub(start).Nanoseconds() / int64(NUM)
	// b.Logf("%dx%d, %d msg(s), %d ns/send, %.1fM msgs/sec", PAR, PAR, NUM, nps, 1.0e03/float64(nps))
	_ = nps
}

func BenchmarkFanIn_Chan_Nx1(b *testing.B) {
	NUM := b.N

	expect := rand.Perm(NUM)
	var count = uint64(NUM)
	var sum = int64(NUM * (NUM - 1) / 2)

	b.ResetTimer()
	start := time.Now()

	channel := NewChanInt(BUFSIZE, 1)

	var rwgbegin, rwgend sync.WaitGroup
	receiver := func() {
		ep, err := channel.NewEndpoint(ReplayAll)
		rwgend.Add(1)
		rwgbegin.Done()
		if err != nil {
			b.Error(err)
			rwgend.Done()
			return
		}
		var rcount uint64
		var rsum int64
		ep.Range(func(value int, err error, closed bool) bool {
			if !closed {
				rsum += int64(value)
				rcount++
			}
			return true
		}, 0)
		if count != rcount {
			b.Errorf("receiver: rcount(%d) != count(%d)", rcount, count /*, "ep=%+v" ep*/)
		} else if sum != rsum {
			b.Errorf("receiver: rsum(%d) != sum(%d)", rsum, sum)
		}
		rwgend.Done()
	}
	rwgbegin.Add(1)
	go receiver()
	rwgbegin.Wait()

	var scount uint64
	var ssum int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := atomic.AddUint64(&scount, 1) - 1
			channel.Send(expect[index])
			atomic.AddInt64(&ssum, int64(expect[index]))
		}
	})
	if sum != ssum {
		b.Errorf("sender: ssum(%d) != sum(%d)", ssum, sum)
	}
	channel.Close(nil)

	rwgend.Wait()

	nps := time.Now().Sub(start).Nanoseconds() / int64(NUM)
	// b.Logf("%dx1, %d msg(s), %d ns/send, %.1fM msgs/sec", runtime.GOMAXPROCS(0), NUM, nps, 1.0e03/float64(nps))
	_ = nps
}

func BenchmarkFanIn_Go_Nx1(b *testing.B) {
	NUM := b.N

	expect := rand.Perm(NUM)
	var count = uint64(NUM)
	var sum = int64(NUM * (NUM - 1) / 2)

	b.ResetTimer()
	start := time.Now()

	c := make(chan int, BUFSIZE)

	wait := make(chan struct{})
	go func() {
		var rcount uint64
		var rsum int64
		for value := range c {
			rsum += int64(value)
			rcount++
		}
		if count != rcount {
			b.Errorf("receiver: rcount(%d) != count(%d)", rcount, count)
		} else if sum != rsum {
			b.Errorf("receiver: rsum(%d) != sum(%d)", rsum, sum)
		}
		close(wait)
	}()

	var scount uint64
	var ssum int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := atomic.AddUint64(&scount, 1) - 1
			c <- expect[index]
			atomic.AddInt64(&ssum, int64(expect[index]))
		}
	})
	if sum != ssum {
		b.Errorf("sender: ssum(%d) != sum(%d)", ssum, sum)
	}

	close(c)
	<-wait

	nps := time.Now().Sub(start).Nanoseconds() / int64(NUM)
	// b.Logf("%dx1, %d msg(s), %d ns/send, %.1fM msgs/sec", runtime.GOMAXPROCS(0), NUM, nps, 1.0e03/float64(nps))
	_ = nps
}

func BenchmarkFanOut_Chan_1xN(b *testing.B) {
	PAR := runtime.GOMAXPROCS(0)
	NUM := b.N

	start := time.Now()

	c := NewChanInt(BUFSIZE, PAR)

	var rwg sync.WaitGroup
	rwg.Add(PAR)

	wait := make(chan struct{})
	go func() {
		rwg.Wait()
		count := 0
		for !c.Closed() {
			c.FastSend(count)
			count++
		}
		close(wait)
	}()

	var rcount int64
	b.RunParallel(func(pb *testing.PB) {
		ep, err := c.NewEndpoint(0)
		rwg.Done()
		if err != nil {
			b.Error(err)
			return
		}
		var sum, count int64
		ep.Range(func(value int, err error, closed bool) bool {
			if !closed && pb.Next() {
				atomic.AddInt64(&rcount, 1)
				sum += int64(value)
				count++
				expectedSum := int64(count) * int64(count-1) / 2
				if sum != expectedSum {
					b.Errorf("data corruption at count == %d ; expected == %d got sum == %d", count, expectedSum, sum)
				}
				return true
			}
			return false
		}, 0)
	})

	c.Close(nil)
	<-wait

	if rcount != int64(NUM) {
		b.Errorf("data loss; expected %d messages got %d", NUM, rcount)
	}

	nps := time.Now().Sub(start).Nanoseconds() / int64(NUM)
	// b.Logf("1x%d, %d msg(s), %d ns/send, %.1fM msgs/sec", PAR, NUM, nps, 1.0e03/float64(nps))
	_ = nps
}

func BenchmarkFanOut_Go_1xN(b *testing.B) {
	PAR := runtime.GOMAXPROCS(0)
	NUMREF := b.N
	NUM := NUMREF / PAR
	if NUM == 0 {
		NUM = 1
	}

	start := time.Now()

	var rwg sync.WaitGroup
	receive := func(ch chan int) {
		var sum, count int64
		for value := range ch {
			sum += int64(value)
			count++
			expectedSum := count * (count - 1) / 2
			if sum != expectedSum {
				b.Errorf("data corruption; at count %d, expected sum %d got %d", count, expectedSum, sum)
			}
		}
		if count != int64(NUM) {
			b.Errorf("data loss; expected %d messages got %d", NUM, count)
		}
		rwg.Done()
	}

	// create channels
	var channels []chan int
	for p := 0; p < PAR; p++ {
		channels = append(channels, make(chan int, BUFSIZE))
	}

	// start receivers
	rwg.Add(PAR)
	for p := 0; p < PAR; p++ {
		go receive(channels[p])
	}

	// send data
	for n := 0; n < NUM; n++ {
		for p := 0; p < PAR; p++ {
			channels[p] <- n
		}
	}

	// close channels
	for p := 0; p < PAR; p++ {
		close(channels[p])
	}

	// wait for receivers
	rwg.Wait()

	nps := time.Now().Sub(start).Nanoseconds() / int64(NUMREF)
	// b.Logf("1x%d, %d msg(s), %d ns/send, %.1fM msgs/sec", PAR, NUMREF, nps, 1.0e03/float64(nps))
	_ = nps
}
