// Package channel provides the Chan type that can multicast and replay messages
// to multiple receivers. It is a specialization of the generics library
// github.com/reactivego/channel/jig on the emtpy interface type.
// 
// The channel can be used by multiple senders to simultaneously send messages
// to the channel that get merged and buffered and are then multicasted to
// multiple concurrent receivers. There is support for a replay buffer to replay
// messages received in the past to newly connecting receivers. Messages are
// timestamped, so during receiving the maximum age of the messages will cause
// old messages to be skipped.
//
//	HOW THIS PACKAGE IS GENERATED 
//
//	The package consists of three files; doc.go, example_test.go and channel.go.
//	The actual package code is in channel.go, but it is generated from a template
//	library using information from doc.go and example_test.go.
//
// 	The [doc.go] file imports the generics library we want to use:
//
//		import _ "github.com/reactivego/channel/jig"
//
// 	The [example_test.go] file contains an example program of all the
// 	functionality that is needed from the generics library. See the Example
//	below this text for the actual source of example_test.go.
//
// 	Generate the package by running the jig command. The code needed to compile
//	the example is generated into the file [channel.go]. That file then contains
//	the code for the actual channel package.
//
//	$ go get github.com/reactivego/channel/jig
//	$ cd $GOPATH/src/github.com/reactivego/channel
//	$ go run -v github.com/reactivego/jig -rv
//	removing file "channel.go"
//	found 16 templates in package "channel" (github.com/reactivego/channel/jig)
//	generating "NewChan"
//	  ChanPadding
//	  ChanState
//	  Chan
//	  ErrOutOfEndpoints
//	  endpoints
//	  NewChan
//	generating "Endpoint"
//	  Endpoint
//	generating "Chan FastSend"
//	  Chan slideBuffer
//	  Chan FastSend
//	generating "Chan Send"
//	  Chan Send
//	generating "Chan Close"
//	  Chan Close
//	generating "Chan Closed"
//	  Chan Closed
//	generating "Chan NewEndpoint"
//	  Chan NewEndpoint
//	generating "Endpoint Range"
//	  Endpoint Range
//	generating "Endpoint Cancel"
//	  Endpoint Cancel
//	generating "Endpoint commitData"
//	  Chan commitData
//	  missing "Endpoint commitData"
//	writing file "channel.go"
package channel

import _ "github.com/reactivego/channel/jig"
