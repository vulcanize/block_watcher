package state_computation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type Trie interface {
	Commit(deleteEmptyObjects bool) (common.Hash, error)
	StateDb() *state.StateDB
}

type StateTrie struct {
	db *state.StateDB
}

func (st *StateTrie) StateDb() *state.StateDB {
	return st.db
}

func (st *StateTrie) Commit(deleteEmptyObjects bool) (common.Hash, error) {
	return st.db.Commit(deleteEmptyObjects)
}
