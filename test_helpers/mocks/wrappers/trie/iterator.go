package trie

import (
	"github.com/ethereum/go-ethereum/common"
)

type MockIterator struct {
	returnHash     common.Hash
	timesToIterate int
}

func NewMockIterator(timesToIterate int) *MockIterator {
	return &MockIterator{
		returnHash:     common.Hash{},
		timesToIterate: timesToIterate,
	}
}

func (mi *MockIterator) SetReturnHash(hash common.Hash) {
	mi.returnHash = hash
}

// TODO: test for path where current node is a leaf (i.e. this returns true)
func (mi *MockIterator) Leaf() bool {
	return false
}

func (mi *MockIterator) LeafBlob() []byte {
	panic("implement me")
}

func (mi *MockIterator) Next(bool) bool {
	if mi.timesToIterate > 0 {
		mi.timesToIterate--
		return true
	}
	return false
}
func (mi *MockIterator) Hash() common.Hash {
	return mi.returnHash
}
