package go_cid

import (
	. "github.com/onsi/gomega"
	"gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"
)

type MockCidDecoder struct {
	passedCidString string
	decodeErr       error
	decodeReturnCid *cid.Cid
}

func NewMockCidDecoder() *MockCidDecoder {
	return &MockCidDecoder{
		passedCidString: "",
		decodeErr:       nil,
		decodeReturnCid: nil,
	}
}

func (decoder *MockCidDecoder) SetDecodeReturnCid(cid *cid.Cid) {
	decoder.decodeReturnCid = cid
}

func (decoder *MockCidDecoder) SetDecodeError(err error) {
	decoder.decodeErr = err
}

func (decoder *MockCidDecoder) Decode(v string) (*cid.Cid, error) {
	decoder.passedCidString = v
	return decoder.decodeReturnCid, decoder.decodeErr
}

func (decoder *MockCidDecoder) AssertDecodeCalledWith(expected string) {
	Expect(decoder.passedCidString).To(Equal(expected))
}
