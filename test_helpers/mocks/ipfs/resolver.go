package ipfs

import (
	. "github.com/onsi/gomega"
	"gx/ipfs/QmTRCUvZLiir12Qr6MV3HKfKMHX8Nf1Vddn6t2g5nsQSb9/go-block-format"
	"gx/ipfs/QmWi2BYBL5gJ3CiAiQchg6rn1A8iBsrWy51EYxvHVjFvLb/go-ipld-format"
)

type MockResolver struct {
	passedBlock blocks.Block
	returnNode  format.Node
	returnErr   error
}

func NewMockResolver() *MockResolver {
	return &MockResolver{
		passedBlock: nil,
		returnNode:  nil,
		returnErr:   nil,
	}
}

func (resolver *MockResolver) SetReturnErr(err error) {
	resolver.returnErr = err
}

func (resolver *MockResolver) SetReturnNode(node format.Node) {
	resolver.returnNode = node
}

func (resolver *MockResolver) Resolve(block blocks.Block) (format.Node, error) {
	resolver.passedBlock = block
	return resolver.returnNode, resolver.returnErr
}

func (resolver *MockResolver) AssertResolveCalledWith(block blocks.Block) {
	Expect(resolver.passedBlock).To(Equal(block))
}
