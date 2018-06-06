package state_mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	. "github.com/onsi/gomega"
)

type MockStateTrie struct {
	commitCalled bool
	returnErr    error
	stateDb      *state.StateDB
}

func NewMockStateTrie() *MockStateTrie {
	return &MockStateTrie{
		commitCalled: false,
		returnErr:    nil,
		stateDb:      nil,
	}
}

func (mst *MockStateTrie) SetReturnErr(err error) {
	mst.returnErr = err
}

func (mst *MockStateTrie) SetStateDb(db *state.StateDB) {
	mst.stateDb = db
}

func (mst *MockStateTrie) Commit(deleteEmptyObjects bool) (common.Hash, error) {
	mst.commitCalled = true
	return common.Hash{}, mst.returnErr
}

func (mst *MockStateTrie) StateDb() *state.StateDB {
	return mst.stateDb
}

func (mst *MockStateTrie) AssertCommitCalled() {
	Expect(mst.commitCalled).To(BeTrue())
}
