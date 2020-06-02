// This file guides regeneration of the heterogeneous multicast package in
// this folder. The [jig tool](https://github.com/reactivego/jig) will generate
// multicast.go guided by the code used in the require function.

// +build ignore

package multicast

import _ "github.com/reactivego/multicast/generic"

func require() {
	c := NewChan(0, 0)
	c.FastSend(nil)
	c.Send(nil)
	c.Close(nil)
	c.Closed()
	e, _ := c.NewEndpoint(ReplayAll)
	e.Range(func(value interface{}, err error, closed bool) bool{ return false }, 0)
	e.Cancel()
}
