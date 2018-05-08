package test_helpers

import . "github.com/onsi/gomega"

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
	bytesToReturn := mp.returnStrings[0]
	mp.returnStrings = mp.returnStrings[1:]
	return bytesToReturn, nil
}

func (mp *MockPublisher) AssertWriteCalledWith(blockDatas [][]byte) {
	Expect(mp.passedBlockDatas).To(Equal(blockDatas))
}
