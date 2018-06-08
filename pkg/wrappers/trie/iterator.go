package trie

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
)

type GethTrieNodeIterator interface {
	Hash() common.Hash
	Leaf() bool
	LeafBlob() []byte
	Next(bool) bool
}

type NodeIterator struct {
	iterator trie.NodeIterator
}

func NewNodeIterator(nodeIterator trie.NodeIterator) *NodeIterator {
	return &NodeIterator{iterator: nodeIterator}
}

func (ni *NodeIterator) Hash() common.Hash {
	return ni.Hash()
}

func (ni *NodeIterator) Leaf() bool {
	return ni.Leaf()
}

func (ni *NodeIterator) LeafBlob() []byte {
	return ni.LeafBlob()
}

func (ni *NodeIterator) Next(b bool) bool {
	return ni.Next(b)
}
