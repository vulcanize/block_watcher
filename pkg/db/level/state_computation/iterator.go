package state_computation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
)

type IIterator interface {
	Hash() common.Hash
	Leaf() bool
	LeafBlob() []byte
	Next(bool) bool
}

type Iterator struct {
	iterator trie.NodeIterator
}

func NewStateIterator(nodeIterator trie.NodeIterator) *Iterator {
	return &Iterator{iterator: nodeIterator}
}

func (si *Iterator) Hash() common.Hash {
	return si.Hash()
}

func (si *Iterator) Leaf() bool {
	return si.Leaf()
}

func (si *Iterator) LeafBlob() []byte {
	return si.LeafBlob()
}

func (si *Iterator) Next(b bool) bool {
	return si.Next(b)
}
