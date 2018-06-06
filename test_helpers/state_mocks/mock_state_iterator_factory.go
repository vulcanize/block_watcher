package state_mocks

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/vulcanize/block_watcher/pkg/db/level/state_computation"
)

type MockStateIteratorFactory struct {
	returnIterator state_computation.Iterator
}

func NewMockStateIteratorFactory() *MockStateIteratorFactory {
	return &MockStateIteratorFactory{}
}

func (msif *MockStateIteratorFactory) SetReturnIterator(iterator state_computation.Iterator) {
	msif.returnIterator = iterator
}

func (msif *MockStateIteratorFactory) NewNodeIterator(stateDB *state.StateDB) state_computation.Iterator {
	return msif.returnIterator
}
