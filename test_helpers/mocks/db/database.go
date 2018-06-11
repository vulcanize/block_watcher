package db

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
)

type MockDatabase struct {
	computeBlockStateTrieErr                          error
	computeBlockStateTriePassedCurrentBlock           *types.Block
	computeBlockStateTriePassedParentBlock            *types.Block
	computeBlockStateTrieReturnHash                   common.Hash
	getBlockBodyByBlockNumberErr                      error
	getBlockBodyByBlockNumberPassedBlockNumbers       []int64
	getBlockBodyByBlockNumberReturnBytes              [][]byte
	getBlockByBlockNumberPassedNumbers                []int64
	getBlockByBlockNumberReturnBlock                  *types.Block
	getBlockHeaderByBlockNumberErr                    error
	getBlockHeaderByBlockNumberPassedBlockNumbers     []int64
	getBlockHeaderByBlockNumberReturnBytes            [][]byte
	getStateAndStorageTrieNodesErr                    error
	getStateAndStorageTrieNodesPassedRoot             common.Hash
	getStateAndStorageTrieNodesReturnStateTrieBytes   [][]byte
	getStateAndStorageTrieNodesReturnStorageTrieBytes [][]byte
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		computeBlockStateTrieErr:                          nil,
		computeBlockStateTriePassedCurrentBlock:           nil,
		computeBlockStateTriePassedParentBlock:            nil,
		computeBlockStateTrieReturnHash:                   common.Hash{},
		getBlockBodyByBlockNumberErr:                      nil,
		getBlockBodyByBlockNumberPassedBlockNumbers:       nil,
		getBlockBodyByBlockNumberReturnBytes:              nil,
		getBlockByBlockNumberPassedNumbers:                nil,
		getBlockByBlockNumberReturnBlock:                  nil,
		getBlockHeaderByBlockNumberErr:                    nil,
		getBlockHeaderByBlockNumberPassedBlockNumbers:     nil,
		getBlockHeaderByBlockNumberReturnBytes:            nil,
		getStateAndStorageTrieNodesErr:                    nil,
		getStateAndStorageTrieNodesPassedRoot:             common.Hash{},
		getStateAndStorageTrieNodesReturnStateTrieBytes:   nil,
		getStateAndStorageTrieNodesReturnStorageTrieBytes: nil,
	}
}

func (md *MockDatabase) SetComputeBlockStateTrieReturnHash(hash common.Hash) {
	md.computeBlockStateTrieReturnHash = hash
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

func (md *MockDatabase) SetGetStateAndStorageTrieNodesReturnStateTrieBytes(returnBytes [][]byte) {
	md.getStateAndStorageTrieNodesReturnStateTrieBytes = returnBytes
}

func (md *MockDatabase) SetGetStateAndStorageTrieNodesReturnStorageTrieBytes(returnBytes [][]byte) {
	md.getStateAndStorageTrieNodesReturnStorageTrieBytes = returnBytes
}

func (md *MockDatabase) SetComputeBlockStateTrieError(err error) {
	md.computeBlockStateTrieErr = err
}

func (md *MockDatabase) SetGetBlockBodyByBlockNumberError(err error) {
	md.getBlockBodyByBlockNumberErr = err
}

func (md *MockDatabase) SetGetBlockHeaderByBlockNumberError(err error) {
	md.getBlockHeaderByBlockNumberErr = err
}

func (md *MockDatabase) SetGetStateTrieNodesError(err error) {
	md.getStateAndStorageTrieNodesErr = err
}

func (md *MockDatabase) ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) (common.Hash, error) {
	md.computeBlockStateTriePassedCurrentBlock = currentBlock
	md.computeBlockStateTriePassedParentBlock = parentBlock
	return md.computeBlockStateTrieReturnHash, md.computeBlockStateTrieErr
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
	md.getBlockByBlockNumberPassedNumbers = append(md.getBlockByBlockNumberPassedNumbers, blockNumber)
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

func (md *MockDatabase) GetStateAndStorageTrieNodes(root common.Hash) ([][]byte, [][]byte, error) {
	md.getStateAndStorageTrieNodesPassedRoot = root
	return md.getStateAndStorageTrieNodesReturnStateTrieBytes, md.getStateAndStorageTrieNodesReturnStorageTrieBytes, md.getStateAndStorageTrieNodesErr
}

func (md *MockDatabase) AssertComputeBlockStateTrieCalledWith(currentBlock *types.Block, parentBlock *types.Block) {
	Expect(md.computeBlockStateTriePassedCurrentBlock).To(Equal(currentBlock))
	Expect(md.computeBlockStateTriePassedParentBlock).To(Equal(parentBlock))
}

func (md *MockDatabase) AssertGetBlockBodyByBlockNumberCalledWith(blockNumbers []int64) {
	Expect(md.getBlockBodyByBlockNumberPassedBlockNumbers).To(Equal(blockNumbers))
}

func (md *MockDatabase) AssertGetBlockByBlockNumberCalledwith(blockNumbers []int64) {
	for i := 0; i < len(blockNumbers); i++ {
		Expect(md.getBlockByBlockNumberPassedNumbers).To(ContainElement(blockNumbers[i]))
	}
	for i := 0; i < len(md.getBlockByBlockNumberPassedNumbers); i++ {
		Expect(blockNumbers).To(ContainElement(md.getBlockByBlockNumberPassedNumbers[i]))
	}
}

func (md *MockDatabase) AssertGetBlockHeaderByBlockNumberCalledWith(blockNumbers []int64) {
	Expect(md.getBlockHeaderByBlockNumberPassedBlockNumbers).To(Equal(blockNumbers))
}

func (md *MockDatabase) AssertGetStateTrieNodesCalledWith(root common.Hash) {
	Expect(md.getStateAndStorageTrieNodesPassedRoot).To(Equal(root))
}
