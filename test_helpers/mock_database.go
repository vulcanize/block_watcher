package test_helpers

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockDatabase struct {
	computeBlockStateTrieErr                      error
	computeBlockStateTriePassedCurrentBlock       *types.Block
	computeBlockStateTriePassedParentBlock        *types.Block
	computeBlockStateTrieReturnBytes              [][]byte
	getBlockBodyByBlockNumberErr                  error
	getBlockBodyByBlockNumberPassedBlockNumbers   []int64
	getBlockBodyByBlockNumberReturnBytes          [][]byte
	getBlockByBlockNumberPassedNumber             int64
	getBlockByBlockNumberReturnBlock              *types.Block
	getBlockHeaderByBlockNumberErr                error
	getBlockHeaderByBlockNumberPassedBlockNumbers []int64
	getBlockHeaderByBlockNumberReturnBytes        [][]byte
	getStateTrieNodesErr                          error
	getStateTrieNodesPassedRoot                   []byte
	getStateTrieNodesReturnBytes                  [][]byte
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		computeBlockStateTrieErr:                      nil,
		computeBlockStateTriePassedCurrentBlock:       nil,
		computeBlockStateTriePassedParentBlock:        nil,
		computeBlockStateTrieReturnBytes:              nil,
		getBlockBodyByBlockNumberErr:                  nil,
		getBlockBodyByBlockNumberPassedBlockNumbers:   nil,
		getBlockBodyByBlockNumberReturnBytes:          nil,
		getBlockByBlockNumberPassedNumber:             0,
		getBlockByBlockNumberReturnBlock:              nil,
		getBlockHeaderByBlockNumberErr:                nil,
		getBlockHeaderByBlockNumberPassedBlockNumbers: nil,
		getBlockHeaderByBlockNumberReturnBytes:        nil,
		getStateTrieNodesErr:                          nil,
		getStateTrieNodesPassedRoot:                   nil,
		getStateTrieNodesReturnBytes:                  nil,
	}
}

func (md *MockDatabase) SetGetBlockBodyByBlockNumberReturnBytes(returnBytes [][]byte) {
	md.getBlockBodyByBlockNumberReturnBytes = returnBytes
}

func (md *MockDatabase) SetGetBlockByBlockNumberReturnBlock(returnBlock *types.Block) {
	md.getBlockByBlockNumberReturnBlock = returnBlock
}

func (md *MockDatabase) SetGetBlockHeaderByBlockNumberReturnBytes(returnBytes [][]byte) {
	md.getBlockHeaderByBlockNumberReturnBytes = returnBytes
}

func (md *MockDatabase) SetGetStateTrieNodesReturnBytes(returnBytes [][]byte) {
	md.getStateTrieNodesReturnBytes = returnBytes
}

func (md *MockDatabase) SetGetBlockBodyByBlockNumberError(err error) {
	md.getBlockBodyByBlockNumberErr = err
}

func (md *MockDatabase) SetGetBlockHeaderByBlockNumberError(err error) {
	md.getBlockHeaderByBlockNumberErr = err
}

func (md *MockDatabase) SetGetStateTrieNodesError(err error) {
	md.getStateTrieNodesErr = err
}

func (md *MockDatabase) ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) ([][]byte, error) {
	md.computeBlockStateTriePassedCurrentBlock = currentBlock
	md.computeBlockStateTriePassedParentBlock = parentBlock
	return md.computeBlockStateTrieReturnBytes, md.computeBlockStateTrieErr
}

func (md *MockDatabase) GetBlockBodyByBlockNumber(blockNumber int64) ([]byte, error) {
	md.getBlockBodyByBlockNumberPassedBlockNumbers = append(md.getBlockBodyByBlockNumberPassedBlockNumbers, blockNumber)
	if md.getBlockBodyByBlockNumberErr != nil {
		return nil, md.getBlockBodyByBlockNumberErr
	}
	returnBytes := md.getBlockBodyByBlockNumberReturnBytes[0]
	md.getBlockBodyByBlockNumberReturnBytes = md.getBlockBodyByBlockNumberReturnBytes[1:]
	return returnBytes, nil
}

func (md *MockDatabase) GetBlockByBlockNumber(blockNumber int64) *types.Block {
	md.getBlockByBlockNumberPassedNumber = blockNumber
	return md.getBlockByBlockNumberReturnBlock
}

func (md *MockDatabase) GetBlockHeaderByBlockNumber(blockNumber int64) ([]byte, error) {
	md.getBlockHeaderByBlockNumberPassedBlockNumbers = append(md.getBlockHeaderByBlockNumberPassedBlockNumbers, blockNumber)
	if md.getBlockHeaderByBlockNumberErr != nil {
		return nil, md.getBlockHeaderByBlockNumberErr
	}
	returnBytes := md.getBlockHeaderByBlockNumberReturnBytes[0]
	md.getBlockHeaderByBlockNumberReturnBytes = md.getBlockHeaderByBlockNumberReturnBytes[1:]
	return returnBytes, nil
}

func (md *MockDatabase) GetStateTrieNodes(root []byte) ([][]byte, error) {
	md.getStateTrieNodesPassedRoot = root
	return md.getStateTrieNodesReturnBytes, md.getStateTrieNodesErr
}

func (md *MockDatabase) AssertComputeBlockStateTrieCalledWith(currentBlock *types.Block, parentBlock *types.Block) {
	Expect(md.computeBlockStateTriePassedCurrentBlock).To(Equal(currentBlock))
	Expect(md.computeBlockStateTriePassedParentBlock).To(Equal(parentBlock))
}

func (md *MockDatabase) AssertGetBlockBodyByBlockNumberCalledWith(blockNumbers []int64) {
	Expect(md.getBlockBodyByBlockNumberPassedBlockNumbers).To(Equal(blockNumbers))
}

func (md *MockDatabase) AssertGetBlockByBlockNumberCalledwith(blockNumber int64) {
	Expect(md.getBlockByBlockNumberPassedNumber).To(Equal(blockNumber))
}

func (md *MockDatabase) AssertGetBlockHeaderByBlockNumberCalledWith(blockNumbers []int64) {
	Expect(md.getBlockHeaderByBlockNumberPassedBlockNumbers).To(Equal(blockNumbers))
}

func (md *MockDatabase) AssertGetStateTrieNodesCalledWith(root []byte) {
	Expect(md.getStateTrieNodesPassedRoot).To(Equal(root))
}
