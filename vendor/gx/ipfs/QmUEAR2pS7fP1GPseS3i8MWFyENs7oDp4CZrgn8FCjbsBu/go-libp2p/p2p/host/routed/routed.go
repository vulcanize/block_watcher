package routedhost

import (
	"context"
	"fmt"
	"time"

	host "gx/ipfs/QmQQGtcp6nVUrQjNsnU53YWV1q8fK1Kd9S7FEkYbRZzxry/go-libp2p-host"

	ma "gx/ipfs/QmUxSEGbv2nmYNnfXi7839wwQqTN3kwQeUxe8dTjZWZs7J/go-multiaddr"
	peer "gx/ipfs/QmVf8hTAsLLFtn4WPCRNdnaF2Eag2qTBS6uR8AiHPZARXy/go-libp2p-peer"
	inet "gx/ipfs/QmXdgNhVEgjLxjUoMs5ViQL7pboAt3Y7V7eGHRiE4qrmTE/go-libp2p-net"
	protocol "gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	pstore "gx/ipfs/QmZhsmorLpD9kmQ4ynbAu4vbKv2goMUnXazwGA4gnWHDjB/go-libp2p-peerstore"
	ifconnmgr "gx/ipfs/Qmav3fJzdn43FDvHyGkPdbQ5JVqqiDPmNdnuGa3vatpmwj/go-libp2p-interface-connmgr"
	msmux "gx/ipfs/QmbXRda5H2K3MSQyWWxTMtd8DWuguEBUCe6hpxfXVpFUGj/go-multistream"
	logging "gx/ipfs/Qmbi1CTJsbnBZjCEgc2otwu8cUFPsGpzWXG7edVCLZ7Gvk/go-log"
	lgbl "gx/ipfs/QmcBbMF4UyZFRTvH9S2h3rbSRBvvEGLqgt4sdvVugG8rX1/go-libp2p-loggables"
)

var log = logging.Logger("routedhost")

// AddressTTL is the expiry time for our addresses.
// We expire them quickly.
const AddressTTL = time.Second * 10

// RoutedHost is a p2p Host that includes a routing system.
// This allows the Host to find the addresses for peers when
// it does not have them.
type RoutedHost struct {
	host  host.Host // embedded other host.
	route Routing
}

type Routing interface {
	FindPeer(context.Context, peer.ID) (pstore.PeerInfo, error)
}

func Wrap(h host.Host, r Routing) *RoutedHost {
	return &RoutedHost{h, r}
}

// Connect ensures there is a connection between this host and the peer with
// given peer.ID. See (host.Host).Connect for more information.
//
// RoutedHost's Connect differs in that if the host has no addresses for a
// given peer, it will use its routing system to try to find some.
func (rh *RoutedHost) Connect(ctx context.Context, pi pstore.PeerInfo) error {
	// first, check if we're already connected.
	if len(rh.Network().ConnsToPeer(pi.ID)) > 0 {
		return nil
	}

	// if we were given some addresses, keep + use them.
	if len(pi.Addrs) > 0 {
		rh.Peerstore().AddAddrs(pi.ID, pi.Addrs, pstore.TempAddrTTL)
	}

	// Check if we have some addresses in our recent memory.
	addrs := rh.Peerstore().Addrs(pi.ID)
	if len(addrs) < 1 {
		// no addrs? find some with the routing system.
		var err error
		addrs, err = rh.findPeerAddrs(ctx, pi.ID)
		if err != nil {
			return err
		}
	}

	// if we're here, we got some addrs. let's use our wrapped host to connect.
	pi.Addrs = addrs
	return rh.host.Connect(ctx, pi)
}

func (rh *RoutedHost) findPeerAddrs(ctx context.Context, id peer.ID) ([]ma.Multiaddr, error) {
	pi, err := rh.route.FindPeer(ctx, id)
	if err != nil {
		return nil, err // couldnt find any :(
	}

	if pi.ID != id {
		err = fmt.Errorf("routing failure: provided addrs for different peer")
		logRoutingErrDifferentPeers(ctx, id, pi.ID, err)
		return nil, err
	}

	return pi.Addrs, nil
}

func logRoutingErrDifferentPeers(ctx context.Context, wanted, got peer.ID, err error) {
	lm := make(lgbl.DeferredMap)
	lm["error"] = err
	lm["wantedPeer"] = func() interface{} { return wanted.Pretty() }
	lm["gotPeer"] = func() interface{} { return got.Pretty() }
	log.Event(ctx, "routingError", lm)
}

func (rh *RoutedHost) ID() peer.ID {
	return rh.host.ID()
}

func (rh *RoutedHost) Peerstore() pstore.Peerstore {
	return rh.host.Peerstore()
}

func (rh *RoutedHost) Addrs() []ma.Multiaddr {
	return rh.host.Addrs()
}

func (rh *RoutedHost) Network() inet.Network {
	return rh.host.Network()
}

func (rh *RoutedHost) Mux() *msmux.MultistreamMuxer {
	return rh.host.Mux()
}

func (rh *RoutedHost) SetStreamHandler(pid protocol.ID, handler inet.StreamHandler) {
	rh.host.SetStreamHandler(pid, handler)
}

func (rh *RoutedHost) SetStreamHandlerMatch(pid protocol.ID, m func(string) bool, handler inet.StreamHandler) {
	rh.host.SetStreamHandlerMatch(pid, m, handler)
}

func (rh *RoutedHost) RemoveStreamHandler(pid protocol.ID) {
	rh.host.RemoveStreamHandler(pid)
}

func (rh *RoutedHost) NewStream(ctx context.Context, p peer.ID, pids ...protocol.ID) (inet.Stream, error) {
	// Ensure we have a connection, with peer addresses resolved by the routing system (#207)
	// It is not sufficient to let the underlying host connect, it will most likely not have
	// any addresses for the peer without any prior connections.
	err := rh.Connect(ctx, pstore.PeerInfo{ID: p})
	if err != nil {
		return nil, err
	}

	return rh.host.NewStream(ctx, p, pids...)
}
func (rh *RoutedHost) Close() error {
	// no need to close IpfsRouting. we dont own it.
	return rh.host.Close()
}
func (rh *RoutedHost) ConnManager() ifconnmgr.ConnManager {
	return rh.host.ConnManager()
}

var _ (host.Host) = (*RoutedHost)(nil)
