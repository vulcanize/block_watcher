package test_helpers

type MockRocksDatabaseReader struct {
	GetBlockHashCalled   bool
	GetBlockHashErr      error
	GetBlockHeaderCalled bool
	GetBlockHeaderErr    error
	OpenDatabaseCalled   bool
	PassedBlockHashKey   []byte
	PassedBlockHeaderKey []byte
	ReturnHashBytes      []byte
	ReturnHeaderBytes    []byte
}

func NewMockRocksDatabaseReader() *MockRocksDatabaseReader {
	return &MockRocksDatabaseReader{
		GetBlockHashCalled:   false,
		GetBlockHashErr:      nil,
		GetBlockHeaderCalled: false,
		GetBlockHeaderErr:    nil,
		OpenDatabaseCalled:   false,
		PassedBlockHashKey:   nil,
		PassedBlockHeaderKey: nil,
		ReturnHashBytes:      nil,
		ReturnHeaderBytes:    nil,
	}
}

func (mrdbr *MockRocksDatabaseReader) SetReturnHashBytes(returnBytes []byte) {
	mrdbr.ReturnHashBytes = returnBytes
}

func (mrdbr *MockRocksDatabaseReader) SetReturnHeaderBytes(returnBytes []byte) {
	mrdbr.ReturnHeaderBytes = returnBytes
}

func (mrdbr *MockRocksDatabaseReader) SetGetHashError(err error) {
	mrdbr.GetBlockHashErr = err
}

func (mrdbr *MockRocksDatabaseReader) SetGetHeaderError(err error) {
	mrdbr.GetBlockHeaderErr = err
}

func (mrdbr *MockRocksDatabaseReader) GetBlockHash(key []byte) ([]byte, error) {
	mrdbr.GetBlockHashCalled = true
	mrdbr.PassedBlockHashKey = key
	return mrdbr.ReturnHashBytes, mrdbr.GetBlockHashErr
}

func (mrdbr *MockRocksDatabaseReader) GetBlockHeader(key []byte) ([]byte, error) {
	mrdbr.GetBlockHeaderCalled = true
	mrdbr.PassedBlockHeaderKey = key
	return mrdbr.ReturnHeaderBytes, mrdbr.GetBlockHeaderErr
}

func (mrdbr *MockRocksDatabaseReader) OpenDatabaseForReadOnlyColumnFamilies(name string) error {
	mrdbr.OpenDatabaseCalled = true
	return nil
}
