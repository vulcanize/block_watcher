package test_helpers

import (
	. "github.com/onsi/gomega"
)

type MockDatabase struct {
	bodyByBlockNumberErr                  error
	bodyByBlockNumberReturnBytes          [][]byte
	bodyByBlockNumberPassedBlockNumbers   []int64
	headerByBlockNumberErr                error
	headerByBlockNumberReturnBytes        [][]byte
	headerByBlockNumberPassedBlockNumbers []int64
	stateTrieNodesErr                     error
	stateTrieNodesPassedRoot              []byte
	stateTrieNodesReturnBytes             [][]byte
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		bodyByBlockNumberErr:                  nil,
		bodyByBlockNumberReturnBytes:          nil,
		bodyByBlockNumberPassedBlockNumbers:   nil,
		headerByBlockNumberErr:                nil,
		headerByBlockNumberReturnBytes:        nil,
		headerByBlockNumberPassedBlockNumbers: nil,
		stateTrieNodesErr:                     nil,
		stateTrieNodesPassedRoot:              nil,
		stateTrieNodesReturnBytes:             nil,
	}
}

func (md *MockDatabase) SetBodyByBlockNumberReturnBytes(returnBytes [][]byte) {
	md.bodyByBlockNumberReturnBytes = returnBytes
}

func (md *MockDatabase) SetHeaderByBlockNumberReturnBytes(returnBytes [][]byte) {
	md.headerByBlockNumberReturnBytes = returnBytes
}

func (md *MockDatabase) SetStateTrieNodesReturnBytes(returnBytes [][]byte) {
	md.stateTrieNodesReturnBytes = returnBytes
}

func (md *MockDatabase) SetBodyByBlockNumberError(err error) {
	md.bodyByBlockNumberErr = err
}

func (md *MockDatabase) SetHeaderByBlockNumberError(err error) {
	md.headerByBlockNumberErr = err
}

func (md *MockDatabase) SetStateTrieNodesError(err error) {
	md.stateTrieNodesErr = err
}

func (md *MockDatabase) GetBlockBodyByBlockNumber(blockNumber int64) ([]byte, error) {
	md.bodyByBlockNumberPassedBlockNumbers = append(md.bodyByBlockNumberPassedBlockNumbers, blockNumber)
	if md.bodyByBlockNumberErr != nil {
		return nil, md.bodyByBlockNumberErr
	}
	returnBytes := md.bodyByBlockNumberReturnBytes[0]
	md.bodyByBlockNumberReturnBytes = md.bodyByBlockNumberReturnBytes[1:]
	return returnBytes, nil
}

func (md *MockDatabase) GetBlockHeaderByBlockNumber(blockNumber int64) ([]byte, error) {
	md.headerByBlockNumberPassedBlockNumbers = append(md.headerByBlockNumberPassedBlockNumbers, blockNumber)
	if md.headerByBlockNumberErr != nil {
		return nil, md.headerByBlockNumberErr
	}
	returnBytes := md.headerByBlockNumberReturnBytes[0]
	md.headerByBlockNumberReturnBytes = md.headerByBlockNumberReturnBytes[1:]
	return returnBytes, nil
}

func (md *MockDatabase) GetStateTrieNodes(root []byte) ([][]byte, error) {
	md.stateTrieNodesPassedRoot = root
	return md.stateTrieNodesReturnBytes, md.stateTrieNodesErr
}

func (md *MockDatabase) AssertGetBlockBodyByBlockNumberCalledWith(blockNumbers []int64) {
	Expect(md.bodyByBlockNumberPassedBlockNumbers).To(Equal(blockNumbers))
}

func (md *MockDatabase) AssertGetBlockHeaderByBlockNumberCalledWith(blockNumbers []int64) {
	Expect(md.headerByBlockNumberPassedBlockNumbers).To(Equal(blockNumbers))
}

func (md *MockDatabase) AssertGetStateTrieNodesCalledWith(root []byte) {
	Expect(md.stateTrieNodesPassedRoot).To(Equal(root))
}
