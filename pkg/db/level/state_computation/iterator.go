package state_computation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type Iterator interface {
	Next() bool
	Hash() common.Hash
}

type StateIterator struct {
	iterator *state.NodeIterator
}

func NewStateIterator(stateNodeIterator *state.NodeIterator) *StateIterator {
	return &StateIterator{iterator: stateNodeIterator}
}

func (si *StateIterator) Next() bool {
	return si.iterator.Next()
}

func (si *StateIterator) Hash() common.Hash {
	return si.iterator.Hash
}
