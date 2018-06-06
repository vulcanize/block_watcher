package state_mocks

import (
	"github.com/ethereum/go-ethereum/common"
)

type MockStateIterator struct {
	returnHash     common.Hash
	timesToIterate int
}

func NewMockStateIterator(timesToIterate int) *MockStateIterator {
	return &MockStateIterator{
		returnHash:     common.Hash{},
		timesToIterate: timesToIterate,
	}
}

func (msi *MockStateIterator) SetReturnHash(hash common.Hash) {
	msi.returnHash = hash
}

func (msi *MockStateIterator) Next() bool {
	if msi.timesToIterate > 0 {
		msi.timesToIterate--
		return true
	}
	return false
}

func (msi *MockStateIterator) Hash() common.Hash {
	return msi.returnHash
}
