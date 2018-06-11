package queue

import peer "gx/ipfs/QmVf8hTAsLLFtn4WPCRNdnaF2Eag2qTBS6uR8AiHPZARXy/go-libp2p-peer"

// PeerQueue maintains a set of peers ordered according to a metric.
// Implementations of PeerQueue could order peers based on distances along
// a KeySpace, latency measurements, trustworthiness, reputation, etc.
type PeerQueue interface {

	// Len returns the number of items in PeerQueue
	Len() int

	// Enqueue adds this node to the queue.
	Enqueue(peer.ID)

	// Dequeue retrieves the highest (smallest int) priority node
	Dequeue() peer.ID
}
