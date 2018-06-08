package state_mocks

import "github.com/vulcanize/block_watcher/pkg/db/level/state_computation"

type MockTrie struct {
	iterator state_computation.IIterator
}

func NewMockTrie() *MockTrie {
	return &MockTrie{}
}

func (mt *MockTrie) SetReturnIterator(iterator state_computation.IIterator) {
	mt.iterator = iterator
}

func (mt *MockTrie) NodeIterator(startKey []byte) state_computation.IIterator {
	return mt.iterator
}
