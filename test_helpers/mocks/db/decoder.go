package db

import . "github.com/onsi/gomega"

type MockDecoder struct {
	called      bool
	err         error
	passedBytes []byte
	passedOut   interface{}
	returnOut   interface{}
}

func NewMockDecoder() *MockDecoder {
	return &MockDecoder{
		called:      false,
		err:         nil,
		passedBytes: nil,
		passedOut:   nil,
		returnOut:   nil,
	}
}

func (md *MockDecoder) SetError(err error) {
	md.err = err
}

func (md *MockDecoder) SetReturnOut(out interface{}) {
	md.returnOut = out
}

func (md *MockDecoder) Decode(raw []byte, out interface{}) (interface{}, error) {
	md.called = true
	md.passedBytes = raw
	md.passedOut = out
	return md.returnOut, md.err
}

func (md *MockDecoder) AssertDecodeCalledWith(raw []byte, out interface{}) {
	Expect(md.called).To(BeTrue())
	Expect(md.passedBytes).To(Equal(raw))
	Expect(md.passedOut).To(BeAssignableToTypeOf(out))
}
