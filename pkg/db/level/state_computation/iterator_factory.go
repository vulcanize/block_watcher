package state_computation

import (
	"github.com/ethereum/go-ethereum/core/state"
)

type IteratorFactory interface {
	NewNodeIterator(stateDB *state.StateDB) Iterator
}

type StateIteratorFactory struct {
}

func NewStateIteratorFactory() *StateIteratorFactory {
	return &StateIteratorFactory{}
}

func (StateIteratorFactory) NewNodeIterator(stateDb *state.StateDB) Iterator {
	stateNodeIterator := state.NewNodeIterator(stateDb)
	return NewStateIterator(stateNodeIterator)
}
