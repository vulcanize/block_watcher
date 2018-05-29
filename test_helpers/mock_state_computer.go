package test_helpers

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockStateComputer struct {
	computeBlockStateTriePassedCurrentBlock *types.Block
	computeBlockStateTriePassedParentBlock  *types.Block
	computeBlockStateTrieReturnErr          error
}

func NewMockStateComputer() *MockStateComputer {
	return &MockStateComputer{
		computeBlockStateTriePassedCurrentBlock: nil,
		computeBlockStateTriePassedParentBlock:  nil,
		computeBlockStateTrieReturnErr:          nil,
	}
}

func (msc *MockStateComputer) SetComputeBlockStateTrieReturnErr(err error) {
	msc.computeBlockStateTrieReturnErr = err
}

func (msc *MockStateComputer) ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) ([][]byte, error) {
	msc.computeBlockStateTriePassedCurrentBlock = currentBlock
	msc.computeBlockStateTriePassedParentBlock = parentBlock
	return nil, msc.computeBlockStateTrieReturnErr
}

func (msc *MockStateComputer) AssertComputeBlockStateTrieCalledWith(currentBlock *types.Block, parentBlock *types.Block) {
	Expect(msc.computeBlockStateTriePassedCurrentBlock).To(Equal(currentBlock))
	Expect(msc.computeBlockStateTriePassedParentBlock).To(Equal(parentBlock))
}
