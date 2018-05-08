package test_helpers

import (
	. "github.com/onsi/gomega"
	ipld "gx/ipfs/Qme5bWv7wtjUNGsK2BNGVUFPKiuxWrsqrtvYwCLRw8YFES/go-ipld-format"
)

type MockAdder struct {
	calledCount int
	passedNodes []ipld.Node
	err         error
}

func NewMockAdder() *MockAdder {
	return &MockAdder{
		calledCount: 0,
		passedNodes: nil,
		err:         nil,
	}
}

func (ma *MockAdder) SetError(err error) {
	ma.err = err
}

func (ma *MockAdder) Add(node ipld.Node) error {
	ma.calledCount++
	ma.passedNodes = append(ma.passedNodes, node)
	return ma.err
}

func (ma *MockAdder) AssertAddCalled(times int, nodeType interface{}) {
	Expect(ma.calledCount).To(Equal(times))
	Expect(len(ma.passedNodes)).To(Equal(times))
	for _, passedNode := range ma.passedNodes {
		Expect(passedNode).To(BeAssignableToTypeOf(nodeType))
	}
}
