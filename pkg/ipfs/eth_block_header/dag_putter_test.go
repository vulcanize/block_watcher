package eth_block_header_test

import (
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/block_watcher/pkg/ipfs/eth_block_header"
	"github.com/vulcanize/block_watcher/test_helpers"
)

var _ = Describe("Creating an IPLD for a block header", func() {
	It("decodes passed bytes into ethereum block header", func() {
		mockDecoder := test_helpers.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		dagPutter := eth_block_header.NewBlockHeaderDagPutter(test_helpers.NewMockAdder(), mockDecoder)
		fakeBytes := []byte{1, 2, 3, 4, 5}

		_, err := dagPutter.DagPut(fakeBytes)

		Expect(err).NotTo(HaveOccurred())
		mockDecoder.AssertDecodeCalledWith(fakeBytes, &types.Header{})
	})

	It("returns error if decoding fails", func() {
		mockDecoder := test_helpers.NewMockDecoder()
		fakeError := errors.New("failed")
		mockDecoder.SetError(fakeError)
		dagPutter := eth_block_header.NewBlockHeaderDagPutter(test_helpers.NewMockAdder(), mockDecoder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakeError))
	})

	It("adds ethereum block header to ipfs", func() {
		mockAdder := test_helpers.NewMockAdder()
		mockDecoder := test_helpers.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		dagPutter := eth_block_header.NewBlockHeaderDagPutter(mockAdder, mockDecoder)
		fakeBytes := []byte{1, 2, 3, 4, 5}

		_, err := dagPutter.DagPut(fakeBytes)

		Expect(err).NotTo(HaveOccurred())
		mockAdder.AssertAddCalled(1, &eth_block_header.EthBlockHeaderNode{})
	})

	It("returns error if adding to ipfs fails", func() {
		mockAdder := test_helpers.NewMockAdder()
		fakeError := errors.New("failed")
		mockAdder.SetError(fakeError)
		mockDecoder := test_helpers.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		dagPutter := eth_block_header.NewBlockHeaderDagPutter(mockAdder, mockDecoder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakeError))
	})
})
