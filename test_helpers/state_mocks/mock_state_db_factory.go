package state_mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/block_watcher/pkg/db/level/state_computation"
)

type MockStateDBFactory struct {
	passedDatabase  state.Database
	passedRoot      common.Hash
	returnErr       error
	returnStateTrie state_computation.IStateDB
}

func NewMockStateDBFactory() *MockStateDBFactory {
	return &MockStateDBFactory{
		passedDatabase:  nil,
		passedRoot:      common.Hash{},
		returnErr:       nil,
		returnStateTrie: nil,
	}
}

func (mstf *MockStateDBFactory) SetStateDB(stateTrie state_computation.IStateDB) {
	mstf.returnStateTrie = stateTrie
}

func (mstf *MockStateDBFactory) SetReturnErr(err error) {
	mstf.returnErr = err
}

func (mstf *MockStateDBFactory) NewStateDB(root common.Hash, db state.Database) (state_computation.IStateDB, error) {
	mstf.passedRoot = root
	mstf.passedDatabase = db
	return mstf.returnStateTrie, mstf.returnErr
}

func (mstf *MockStateDBFactory) AssertNewStateTrieCalledWith(root common.Hash, db state.Database) {
	Expect(mstf.passedRoot).To(Equal(root))
	Expect(mstf.passedDatabase).To(Equal(db))
}
