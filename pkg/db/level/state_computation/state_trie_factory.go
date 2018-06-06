package state_computation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type TrieFactory interface {
	NewStateTrie(root common.Hash, db state.Database) (Trie, error)
}

type StateTrieFactory struct {
}

func NewStateTrieFactory() *StateTrieFactory {
	return &StateTrieFactory{}
}

func (stf *StateTrieFactory) NewStateTrie(root common.Hash, db state.Database) (Trie, error) {
	stateDb, err := state.New(root, db)
	if err != nil {
		return nil, err
	}
	return &StateTrie{db: stateDb}, nil
}
