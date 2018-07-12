package ipfs

import (
	. "github.com/onsi/gomega"
	"gx/ipfs/QmTRCUvZLiir12Qr6MV3HKfKMHX8Nf1Vddn6t2g5nsQSb9/go-block-format"
	"gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"
)

type MockGetter struct {
	passedCid   *cid.Cid
	returnBlock blocks.Block
	returnErr   error
}

func NewMockGetter() *MockGetter {
	return &MockGetter{
		passedCid:   nil,
		returnBlock: nil,
		returnErr:   nil,
	}
}

func (getter *MockGetter) SetReturnBlock(block blocks.Block) {
	getter.returnBlock = block
}

func (getter *MockGetter) SetReturnErr(err error) {
	getter.returnErr = err
}

func (getter *MockGetter) Get(cid *cid.Cid) (blocks.Block, error) {
	getter.passedCid = cid
	return getter.returnBlock, getter.returnErr
}

func (getter *MockGetter) AssertGetCalledWith(cid *cid.Cid) {
	Expect(getter.passedCid).To(Equal(cid))
}
