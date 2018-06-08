package state_mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	. "github.com/onsi/gomega"
)

type MockStateDB struct {
	commitCalled bool
	returnErr    error
	stateDb      *state.StateDB
}

func NewMockStateDB() *MockStateDB {
	return &MockStateDB{
		commitCalled: false,
		returnErr:    nil,
		stateDb:      nil,
	}
}

func (mst *MockStateDB) SetReturnErr(err error) {
	mst.returnErr = err
}

func (mst *MockStateDB) SetStateDB(db *state.StateDB) {
	mst.stateDb = db
}

func (mst *MockStateDB) Commit(deleteEmptyObjects bool) (common.Hash, error) {
	mst.commitCalled = true
	return common.Hash{}, mst.returnErr
}

func (mst *MockStateDB) StateDB() *state.StateDB {
	return mst.stateDb
}

func (mst *MockStateDB) AssertCommitCalled() {
	Expect(mst.commitCalled).To(BeTrue())
}
