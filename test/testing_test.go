package test

import (
	"runtime"
	"testing"
	"time"
)

func TestSleepingReceiver(t *testing.T) {
	channel := NewChanInt(128, 1)
	ep, err := channel.NewEndpoint(ReplayAll)
	if err != nil {
		t.Error(err)
	}
	wait := make(chan struct{})
	go func() {
		ep.Range(func(value int, err error, closed bool) bool {
			if !closed {
			}
			return true
		}, 0)
		close(wait)
	}()
	time.Sleep(300 * time.Millisecond)
	channel.Send(1)
	channel.Close(nil)
	<-wait
}

func TestChanMaxAge(t *testing.T) {
	channel := NewChanInt(128, 1)
	ep, err := channel.NewEndpoint(ReplayAll)
	if err != nil {
		t.Error(err)
	}

	start := time.Now()
	for i := 0; i < 100; i++ {
		// wait until next millisecond.
		for time.Since(start) < time.Duration(i)*time.Millisecond {
			runtime.Gosched()
		}
		channel.Send(i)
	}
	channel.Close(nil)

	num := 50
	count := func(value int, err error, closed bool) bool {
		if !closed {
			if value != num {
				t.Errorf("expected %d, got %d", num, value)
			}
			num++
		}
		return true
	}
	ep.Range(count /*49.5ms*/, 99*(time.Millisecond/2))
}

func TestChanEndpointKeep(t *testing.T) {
	channel := NewChanInt(128, 1)
	for i := 0; i < 100; i++ {
		channel.Send(i)
	}
	channel.Close(nil)

	ep, err := channel.NewEndpoint(0)
	if err != nil {
		t.Fatal(err)
	}
	num := 0
	ep.Range(func(value int, err error, closed bool) bool {
		if !closed {
			num++
		}
		return true
	}, 0)
	if num != 0 {
		t.Fatal("Got", num, "buffered values but I ask for none (keep arg was 0)")
	}
}
