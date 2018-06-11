package relay

import (
	"context"
	"fmt"
	"math/rand"

	ma "gx/ipfs/QmUxSEGbv2nmYNnfXi7839wwQqTN3kwQeUxe8dTjZWZs7J/go-multiaddr"
	tpt "gx/ipfs/QmVMBFZqRZDA6TrQkVGGJEDSp5jC3UUMUjLcvaZ3fLCqh4/go-libp2p-transport"
	peer "gx/ipfs/QmVf8hTAsLLFtn4WPCRNdnaF2Eag2qTBS6uR8AiHPZARXy/go-libp2p-peer"
	pstore "gx/ipfs/QmZhsmorLpD9kmQ4ynbAu4vbKv2goMUnXazwGA4gnWHDjB/go-libp2p-peerstore"
)

func (d *RelayTransport) Dial(ctx context.Context, a ma.Multiaddr, p peer.ID) (tpt.Conn, error) {
	c, err := d.Relay().Dial(ctx, a)
	if err != nil {
		return nil, err
	}
	return d.upgrader.UpgradeOutbound(ctx, d, c, p)
}

func (r *Relay) Dial(ctx context.Context, a ma.Multiaddr) (*Conn, error) {
	if !r.Matches(a) {
		return nil, fmt.Errorf("%s is not a relay address", a)
	}
	parts := ma.Split(a)

	spl, _ := ma.NewMultiaddr("/p2p-circuit")

	var relayaddr, destaddr ma.Multiaddr
	for i, p := range parts {
		if p.Equal(spl) {
			relayaddr = ma.Join(parts[:i]...)
			destaddr = ma.Join(parts[i+1:]...)
			break
		}
	}

	dinfo, err := pstore.InfoFromP2pAddr(destaddr)
	if err != nil {
		return nil, err
	}

	if len(relayaddr.Bytes()) == 0 {
		// unspecific relay address, try dialing using known hop relays
		return r.tryDialRelays(ctx, *dinfo)
	}

	var rinfo *pstore.PeerInfo
	rinfo, err = pstore.InfoFromP2pAddr(relayaddr)
	if err != nil {
		return nil, err
	}

	return r.DialPeer(ctx, *rinfo, *dinfo)
}

func (r *Relay) tryDialRelays(ctx context.Context, dinfo pstore.PeerInfo) (*Conn, error) {
	var relays []peer.ID
	r.mx.Lock()
	for p := range r.relays {
		relays = append(relays, p)
	}
	r.mx.Unlock()

	// shuffle list of relays, avoid overloading a specific relay
	for i := range relays {
		j := rand.Intn(i + 1)
		relays[i], relays[j] = relays[j], relays[i]
	}

	for _, relay := range relays {
		if len(r.host.Network().ConnsToPeer(relay)) == 0 {
			continue
		}

		rctx, cancel := context.WithTimeout(ctx, HopConnectTimeout)
		c, err := r.DialPeer(rctx, pstore.PeerInfo{ID: relay}, dinfo)
		cancel()

		if err == nil {
			return c, nil
		}

		log.Debugf("error opening relay connection through %s: %s", dinfo.ID, err.Error())
	}

	return nil, fmt.Errorf("Failed to dial through %d known relay hosts", len(relays))
}
