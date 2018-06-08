package state_computation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type IStateDBFactory interface {
	NewStateDB(root common.Hash, db state.Database) (IStateDB, error)
}

type StateDBFactory struct {
}

func NewStateTrieFactory() *StateDBFactory {
	return &StateDBFactory{}
}

func (stf *StateDBFactory) NewStateDB(root common.Hash, db state.Database) (IStateDB, error) {
	stateDb, err := state.New(root, db)
	if err != nil {
		return nil, err
	}
	return &StateDB{db: stateDb}, nil
}
