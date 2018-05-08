package test_helpers

type MockDecompressor struct {
	Called      bool
	Err         error
	PassedBytes []byte
}

func NewMockDecompressor() *MockDecompressor {
	return &MockDecompressor{
		Called:      false,
		Err:         nil,
		PassedBytes: nil,
	}
}

func (md *MockDecompressor) SetError(err error) {
	md.Err = err
}

func (md *MockDecompressor) Decompress(raw []byte) ([]byte, error) {
	md.Called = true
	md.PassedBytes = raw
	return nil, md.Err
}
