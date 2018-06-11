package ipfs

import (
	. "github.com/onsi/gomega"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
)

type MockPublisher struct {
	err              error
	passedBlockDatas [][]byte
	returnStrings    [][]string
}

func NewMockPublisher() *MockPublisher {
	return &MockPublisher{
		err:              nil,
		passedBlockDatas: [][]byte{},
		returnStrings:    nil,
	}
}

func (mp *MockPublisher) SetReturnStrings(returnBytes [][]string) {
	mp.returnStrings = returnBytes
}

func (mp *MockPublisher) SetError(err error) {
	mp.err = err
}

func (mp *MockPublisher) Write(blockData []byte) ([]string, error) {
	mp.passedBlockDatas = append(mp.passedBlockDatas, blockData)
	if mp.err != nil {
		return nil, mp.err
	}
	var stringsToReturn []string
	if len(mp.returnStrings) > 0 {
		stringsToReturn = mp.returnStrings[0]
		if len(mp.returnStrings) > 1 {
			mp.returnStrings = mp.returnStrings[1:]
		} else {
			mp.returnStrings = [][]string{{test_helpers.FakeString}}
		}
	} else {
		stringsToReturn = []string{test_helpers.FakeString}
	}
	return stringsToReturn, nil
}

func (mp *MockPublisher) AssertWriteCalledWith(blockDatas [][]byte) {
	for i := 0; i < len(blockDatas); i++ {
		Expect(mp.passedBlockDatas).To(ContainElement(blockDatas[i]))
	}
	for i := 0; i < len(mp.passedBlockDatas); i++ {
		Expect(blockDatas).To(ContainElement(mp.passedBlockDatas[i]))
	}
}
