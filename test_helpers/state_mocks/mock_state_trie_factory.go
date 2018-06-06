package state_mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/block_watcher/pkg/db/level/state_computation"
)

type MockStateTrieFactory struct {
	passedDatabase  state.Database
	passedRoot      common.Hash
	returnErr       error
	returnStateTrie state_computation.Trie
}

func NewMockStateTrieFactory() *MockStateTrieFactory {
	return &MockStateTrieFactory{
		passedDatabase:  nil,
		passedRoot:      common.Hash{},
		returnErr:       nil,
		returnStateTrie: nil,
	}
}

func (mstf *MockStateTrieFactory) SetStateTrie(stateTrie state_computation.Trie) {
	mstf.returnStateTrie = stateTrie
}

func (mstf *MockStateTrieFactory) SetReturnErr(err error) {
	mstf.returnErr = err
}

func (mstf *MockStateTrieFactory) NewStateTrie(root common.Hash, db state.Database) (state_computation.Trie, error) {
	mstf.passedRoot = root
	mstf.passedDatabase = db
	return mstf.returnStateTrie, mstf.returnErr
}

func (mstf *MockStateTrieFactory) AssertNewStateTrieCalledWith(root common.Hash, db state.Database) {
	Expect(mstf.passedRoot).To(Equal(root))
	Expect(mstf.passedDatabase).To(Equal(db))
}
