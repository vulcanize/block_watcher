// Package mockrouting provides a virtual routing server. To use it,
// create a virtual routing server and use the Client() method to get a
// routing client (IpfsRouting). The server quacks like a DHT but is
// really a local in-memory hash table.
package mockrouting

import (
	"context"

	"gx/ipfs/QmPdxCaVp4jZ9RbxqZADvKH6kiCR5jHvdR5f2ycjAY6T2a/go-testutil"
	delay "gx/ipfs/QmRJVNatYJwTAHgdSM1Xef9QVQ1Ch3XHdmcrykjP5Y4soL/go-ipfs-delay"
	routing "gx/ipfs/QmUV9hDAAyjeGbxbXkJ2sYqZ6dTd1DXJ2REhYEkRm178Tg/go-libp2p-routing"
	peer "gx/ipfs/QmVf8hTAsLLFtn4WPCRNdnaF2Eag2qTBS6uR8AiHPZARXy/go-libp2p-peer"
	ds "gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore"
)

// Server provides mockrouting Clients
type Server interface {
	Client(p testutil.Identity) Client
	ClientWithDatastore(context.Context, testutil.Identity, ds.Datastore) Client
}

// Client implements IpfsRouting
type Client interface {
	routing.IpfsRouting
}

// NewServer returns a mockrouting Server
func NewServer() Server {
	return NewServerWithDelay(DelayConfig{
		ValueVisibility: delay.Fixed(0),
		Query:           delay.Fixed(0),
	})
}

// NewServerWithDelay returns a mockrouting Server with a delay!
func NewServerWithDelay(conf DelayConfig) Server {
	return &s{
		providers: make(map[string]map[peer.ID]providerRecord),
		delayConf: conf,
	}
}

// DelayConfig can be used to configured the fake delays of a mock server.
// Use with NewServerWithDelay().
type DelayConfig struct {
	// ValueVisibility is the time it takes for a value to be visible in the network
	// FIXME there _must_ be a better term for this
	ValueVisibility delay.D

	// Query is the time it takes to receive a response from a routing query
	Query delay.D
}
