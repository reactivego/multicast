// go test -run=integrity -parallel=10 -cpu=1,8,10

package test

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChan_FanInOut_integrity(t *testing.T) {
	var numSenders = runtime.GOMAXPROCS(0)
	var numReceivers = runtime.GOMAXPROCS(0)

	const permutations = 0x200000
	expect := rand.Perm(permutations)

	channel := NewChanInt(BUFSIZE, numReceivers)

	var count uint64 = permutations
	var sum int64
	for i := 0; i < permutations; i++ {
		sum += int64(expect[i])
	}

	var numgoroutines uint32
	var rwg sync.WaitGroup

	var scount uint64
	var ssum int64
	var swg sync.WaitGroup
	sender := func(name string) {
		atomic.AddUint32(&numgoroutines, 1)

		// println(name)
		swg.Add(1)
		rwg.Wait()

		index := atomic.AddUint64(&scount, 1) - 1
		for index < permutations {
			channel.Send(expect[index])
			atomic.AddInt64(&ssum, int64(expect[index]))
			index = atomic.AddUint64(&scount, 1) - 1
		}

		swg.Done()
		swg.Wait() // wait for all senders to complete.
		channel.Close(nil)
	}

	receiver := func(t *testing.T) {
		rwg.Add(1)
		ep, err := channel.NewEndpoint(ReplayAll)
		if err != nil {
			assert.NoError(t, err)
			channel.Close(nil)
			return
		}
		rwg.Done()

		t.Parallel()
		atomic.AddUint32(&numgoroutines, 1)

		// println(t.Name())
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
			assert.Equalf(t, count, rcount, "ep=%+v", ep)
			return
		}
		assert.Equal(t, sum, rsum)
	}

	rwg.Add(1)
	for i := 0; i < numSenders; i++ {
		name := fmt.Sprintf("Sender%d", i)
		go sender(name)
	}

	t.Run("receivers", func(t *testing.T) {
		for i := 0; i < numReceivers; i++ {
			t.Run(fmt.Sprintf("Receiver%d", i), receiver)
		}
		rwg.Done()
	})

	assert.Equal(t, sum, ssum)
}
