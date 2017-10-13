package test

import (
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestMemoryLayout(t *testing.T) {
	const sizeofTime = 24
	var tm time.Time
	result := int(unsafe.Sizeof(tm))
	assert.Equal(t, sizeofTime, result)

	const sizeoferror = 16
	var e error
	result = int(unsafe.Sizeof(e))
	assert.Equal(t, sizeoferror, result)

	const sizeofinterface = 16
	var i interface{}
	result = int(unsafe.Sizeof(i))
	assert.Equal(t, sizeofinterface, result)

	const sizeofSlice = 24
	var s []interface{}
	result = int(unsafe.Sizeof(s))
	assert.Equal(t, sizeofSlice, result)

	const sizeofPointer = 8
	p := &ChanInt{}
	result = int(unsafe.Sizeof(p))
	assert.Equal(t, sizeofPointer, result)

	type test struct {
		first  []struct{}
		second uint32
		pad0   [0]byte // removed
		third  uint32
		pad1   [0]byte // forces 8 byte padding
	}
	const sizeofStruct = sizeofSlice + 4 + 4 + 8
	st := test{}
	result = int(unsafe.Sizeof(st))
	assert.Equal(t, sizeofStruct, result)
	assert.Equal(t, sizeofSlice, int(unsafe.Offsetof(st.second)))
	assert.Equal(t, sizeofSlice+4, int(unsafe.Offsetof(st.third)))
	assert.Equal(t, sizeofSlice+4+4, int(unsafe.Offsetof(st.pad1)))

	type endpoint struct {
		*int
		_____________a pad56
		cursor         uint64
		_____________b pad56
		endpointState  uint64
		_____________c pad56
		lastActive     time.Time
		_____________d pad40
		endpointClosed uint64
		_____________e pad56
	}

	const sizeofendpoints = _PADDING*(_EXTRA_PADDING+(24+4+4+(32))) + (1-_PADDING)*(24+4+4+(8))
	eps := struct {
		entry    []endpoint
		len      uint32
		activity uint32 // idling, enumerating, creating
		________ pad32
	}{}
	result = int(unsafe.Sizeof(eps))
	assert.Equal(t, sizeofendpoints, result)

	const sizeofEndpoint = _PADDING*(5*(_EXTRA_PADDING+64)) + (1-_PADDING)*(16+8+8+24+8)
	ep := endpoint{}
	result = int(unsafe.Sizeof(ep))
	assert.Equal(t, sizeofEndpoint, result)
}
