package test_helpers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	. "github.com/onsi/gomega"
)

type MockLevelDatabaseReader struct {
	getBodyRLPPassedHash         common.Hash
	getBodyRLPPassedNumber       uint64
	getCanonicalHashPassedNumber uint64
	getCanonicalHashReturnHash   common.Hash
	getHeaderRLPPassedHash       common.Hash
	getHeaderRLPPassedNumber     uint64
	getStateTrieNodesPassedRoot  common.Hash
	getStateTrieNodesReturnErr   error
	getStateTrieNodesReturnBytes [][]byte
}

func NewMockLevelDatabaseReader() *MockLevelDatabaseReader {
	return &MockLevelDatabaseReader{
		getBodyRLPPassedHash:         common.Hash{},
		getBodyRLPPassedNumber:       0,
		getCanonicalHashPassedNumber: 0,
		getCanonicalHashReturnHash:   common.Hash{},
		getHeaderRLPPassedHash:       common.Hash{},
		getHeaderRLPPassedNumber:     0,
		getStateTrieNodesPassedRoot:  common.Hash{},
		getStateTrieNodesReturnErr:   nil,
		getStateTrieNodesReturnBytes: nil,
	}
}

func (mldr *MockLevelDatabaseReader) SetGetCanonicalHashReturnHash(hash common.Hash) {
	mldr.getCanonicalHashReturnHash = hash
}

func (mldr *MockLevelDatabaseReader) SetGetStateTrieNodesReturnBytes(returnBytes [][]byte) {
	mldr.getStateTrieNodesReturnBytes = returnBytes
}

func (mldr *MockLevelDatabaseReader) SetGetStateTrieNodesReturnErr(err error) {
	mldr.getStateTrieNodesReturnErr = err
}

func (mldr *MockLevelDatabaseReader) GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue {
	mldr.getBodyRLPPassedHash = hash
	mldr.getBodyRLPPassedNumber = number
	return nil
}

func (mldr *MockLevelDatabaseReader) GetCanonicalHash(number uint64) common.Hash {
	mldr.getCanonicalHashPassedNumber = number
	return mldr.getCanonicalHashReturnHash
}

func (mldr *MockLevelDatabaseReader) GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	mldr.getHeaderRLPPassedHash = hash
	mldr.getHeaderRLPPassedNumber = number
	return nil
}

func (mldr *MockLevelDatabaseReader) GetStateTrieNodes(root common.Hash) ([][]byte, error) {
	mldr.getStateTrieNodesPassedRoot = root
	return mldr.getStateTrieNodesReturnBytes, mldr.getStateTrieNodesReturnErr
}

func (mldr *MockLevelDatabaseReader) AssertGetBodyRLPCalledWith(hash common.Hash, number uint64) {
	Expect(mldr.getBodyRLPPassedHash).To(Equal(hash))
	Expect(mldr.getBodyRLPPassedNumber).To(Equal(number))
}

func (mldr *MockLevelDatabaseReader) AssertGetCanonicalHashCalledWith(number uint64) {
	Expect(mldr.getCanonicalHashPassedNumber).To(Equal(number))
}

func (mldr *MockLevelDatabaseReader) AssertGetHeaderRLPCalledWith(hash common.Hash, number uint64) {
	Expect(mldr.getHeaderRLPPassedHash).To(Equal(hash))
	Expect(mldr.getHeaderRLPPassedNumber).To(Equal(number))
}

func (mldr *MockLevelDatabaseReader) AssertGetStateTrieNodesCalledWith(root common.Hash) {
	Expect(mldr.getStateTrieNodesPassedRoot).To(Equal(root))
}
