package rawdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	. "github.com/onsi/gomega"
)

type MockAccessorsChain struct {
	getBlockPassedHash                                common.Hash
	getBlockPassedNumber                              uint64
	getBlockReturnBlock                               *types.Block
	getBodyRLPPassedHash                              common.Hash
	getBodyRLPPassedNumber                            uint64
	getCanonicalHashPassedNumber                      uint64
	getCanonicalHashReturnHash                        common.Hash
	getHeaderRLPPassedHash                            common.Hash
	getHeaderRLPPassedNumber                          uint64
	getStateAndStorageTrieNodesPassedRoot             common.Hash
	getStateAndStorageTrieNodesReturnErr              error
	getStateAndStorageTrieNodesReturnStateTrieBytes   [][]byte
	getStateAndStorageTrieNodesReturnStorageTrieBytes [][]byte
}

func NewMockAccessorsChain() *MockAccessorsChain {
	return &MockAccessorsChain{
		getBlockPassedHash:                                common.Hash{},
		getBlockPassedNumber:                              0,
		getBlockReturnBlock:                               nil,
		getBodyRLPPassedHash:                              common.Hash{},
		getBodyRLPPassedNumber:                            0,
		getCanonicalHashPassedNumber:                      0,
		getCanonicalHashReturnHash:                        common.Hash{},
		getHeaderRLPPassedHash:                            common.Hash{},
		getHeaderRLPPassedNumber:                          0,
		getStateAndStorageTrieNodesPassedRoot:             common.Hash{},
		getStateAndStorageTrieNodesReturnErr:              nil,
		getStateAndStorageTrieNodesReturnStateTrieBytes:   nil,
		getStateAndStorageTrieNodesReturnStorageTrieBytes: nil,
	}
}

func (mldr *MockAccessorsChain) SetGetBlockReturnBlock(returnBlock *types.Block) {
	mldr.getBlockReturnBlock = returnBlock
}

func (mldr *MockAccessorsChain) SetGetCanonicalHashReturnHash(hash common.Hash) {
	mldr.getCanonicalHashReturnHash = hash
}

func (mldr *MockAccessorsChain) SetGetStateTrieNodesReturnStateTrieBytes(returnBytes [][]byte) {
	mldr.getStateAndStorageTrieNodesReturnStateTrieBytes = returnBytes
}

func (mldr *MockAccessorsChain) SetGetStateTrieNodesReturnStorageTrieBytes(returnBytes [][]byte) {
	mldr.getStateAndStorageTrieNodesReturnStorageTrieBytes = returnBytes
}

func (mldr *MockAccessorsChain) SetGetStateTrieNodesReturnErr(err error) {
	mldr.getStateAndStorageTrieNodesReturnErr = err
}

func (mldr *MockAccessorsChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	mldr.getBlockPassedHash = hash
	mldr.getBlockPassedNumber = number
	return mldr.getBlockReturnBlock
}

func (mldr *MockAccessorsChain) GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue {
	mldr.getBodyRLPPassedHash = hash
	mldr.getBodyRLPPassedNumber = number
	return nil
}

func (mldr *MockAccessorsChain) GetCanonicalHash(number uint64) common.Hash {
	mldr.getCanonicalHashPassedNumber = number
	return mldr.getCanonicalHashReturnHash
}

func (mldr *MockAccessorsChain) GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	mldr.getHeaderRLPPassedHash = hash
	mldr.getHeaderRLPPassedNumber = number
	return nil
}

func (mldr *MockAccessorsChain) GetStateAndStorageTrieNodes(root common.Hash) ([][]byte, [][]byte, error) {
	mldr.getStateAndStorageTrieNodesPassedRoot = root
	return mldr.getStateAndStorageTrieNodesReturnStateTrieBytes, mldr.getStateAndStorageTrieNodesReturnStorageTrieBytes, mldr.getStateAndStorageTrieNodesReturnErr
}

func (mldr *MockAccessorsChain) AssertGetBlockCalledWith(hash common.Hash, number uint64) {
	Expect(mldr.getBlockPassedHash).To(Equal(hash))
	Expect(mldr.getBlockPassedNumber).To(Equal(number))
}

func (mldr *MockAccessorsChain) AssertGetBodyRLPCalledWith(hash common.Hash, number uint64) {
	Expect(mldr.getBodyRLPPassedHash).To(Equal(hash))
	Expect(mldr.getBodyRLPPassedNumber).To(Equal(number))
}

func (mldr *MockAccessorsChain) AssertGetCanonicalHashCalledWith(number uint64) {
	Expect(mldr.getCanonicalHashPassedNumber).To(Equal(number))
}

func (mldr *MockAccessorsChain) AssertGetHeaderRLPCalledWith(hash common.Hash, number uint64) {
	Expect(mldr.getHeaderRLPPassedHash).To(Equal(hash))
	Expect(mldr.getHeaderRLPPassedNumber).To(Equal(number))
}

func (mldr *MockAccessorsChain) AssertGetStateTrieNodesCalledWith(root common.Hash) {
	Expect(mldr.getStateAndStorageTrieNodesPassedRoot).To(Equal(root))
}
