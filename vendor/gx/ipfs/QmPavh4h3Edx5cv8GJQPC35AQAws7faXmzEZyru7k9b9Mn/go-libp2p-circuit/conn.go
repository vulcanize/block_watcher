package relay

import (
	"fmt"
	"net"

	manet "gx/ipfs/QmNqRnejxJxjRroz7buhrjfU8i3yNBLa81hFtmf2pXEffN/go-multiaddr-net"
	ma "gx/ipfs/QmUxSEGbv2nmYNnfXi7839wwQqTN3kwQeUxe8dTjZWZs7J/go-multiaddr"
	inet "gx/ipfs/QmXdgNhVEgjLxjUoMs5ViQL7pboAt3Y7V7eGHRiE4qrmTE/go-libp2p-net"
	pstore "gx/ipfs/QmZhsmorLpD9kmQ4ynbAu4vbKv2goMUnXazwGA4gnWHDjB/go-libp2p-peerstore"
)

type Conn struct {
	inet.Stream
	remote pstore.PeerInfo
}

type NetAddr struct {
	Relay  string
	Remote string
}

func (n *NetAddr) Network() string {
	return "libp2p-circuit-relay"
}

func (n *NetAddr) String() string {
	return fmt.Sprintf("relay[%s-%s]", n.Remote, n.Relay)
}

func (c *Conn) RemoteAddr() net.Addr {
	return &NetAddr{
		Relay:  c.Conn().RemotePeer().Pretty(),
		Remote: c.remote.ID.Pretty(),
	}
}

func (c *Conn) RemoteMultiaddr() ma.Multiaddr {
	a, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s/p2p-circuit/ipfs/%s", c.Conn().RemotePeer().Pretty(), c.remote.ID.Pretty()))
	if err != nil {
		panic(err)
	}
	return a
}

func (c *Conn) LocalMultiaddr() ma.Multiaddr {
	return c.Conn().LocalMultiaddr()
}

func (c *Conn) LocalAddr() net.Addr {
	na, err := manet.ToNetAddr(c.Conn().LocalMultiaddr())
	if err != nil {
		log.Error("failed to convert local multiaddr to net addr:", err)
		return nil
	}
	return na
}
