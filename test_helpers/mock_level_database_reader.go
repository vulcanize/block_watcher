package test_helpers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type MockLevelDatabaseReader struct {
	GetBodyCalled   bool
	GetHashCalled   bool
	GetHeaderCalled bool
	PassedHash      common.Hash
	PassedNumber    uint64
	ReturnHash      common.Hash
}

func NewMockLevelDatabaseReader() *MockLevelDatabaseReader {
	return &MockLevelDatabaseReader{
		GetBodyCalled:   false,
		GetHashCalled:   false,
		GetHeaderCalled: false,
		PassedHash:      common.Hash{},
		PassedNumber:    0,
		ReturnHash:      common.Hash{},
	}
}

func (mldr *MockLevelDatabaseReader) SetReturnHash(hash common.Hash) {
	mldr.ReturnHash = hash
}

func (mldr *MockLevelDatabaseReader) GetCanonicalHash(number uint64) common.Hash {
	mldr.GetHashCalled = true
	mldr.PassedNumber = number
	return mldr.ReturnHash
}

func (mldr *MockLevelDatabaseReader) GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	mldr.GetHeaderCalled = true
	mldr.PassedHash = hash
	mldr.PassedNumber = number
	return nil
}

func (mldr *MockLevelDatabaseReader) GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue {
	mldr.GetBodyCalled = true
	mldr.PassedHash = hash
	mldr.PassedNumber = number
	return nil
}
