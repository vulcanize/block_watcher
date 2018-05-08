package test_helpers

type MockDagPutter struct {
	Called      bool
	PassedBytes []byte
	Err         error
}

func NewMockDagPutter() *MockDagPutter {
	return &MockDagPutter{
		Called:      false,
		PassedBytes: nil,
		Err:         nil,
	}
}

func (mdp *MockDagPutter) SetError(err error) {
	mdp.Err = err
}

func (mdp *MockDagPutter) DagPut(raw []byte) ([]string, error) {
	mdp.Called = true
	mdp.PassedBytes = raw
	return nil, mdp.Err
}
