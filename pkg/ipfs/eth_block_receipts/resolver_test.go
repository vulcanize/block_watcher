package eth_block_receipts_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gx/ipfs/QmTRCUvZLiir12Qr6MV3HKfKMHX8Nf1Vddn6t2g5nsQSb9/go-block-format"

	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_receipts"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/go-ethereum/rlp"
)

var _ = Describe("Eth transaction receipt resolver", func() {
	It("decodes block's raw data into receipt", func() {
		block := blocks.NewBlock(test_helpers.FakeReceiptRLP)
		mockDecoder := rlp.NewMockDecoder()
		resolver := eth_block_receipts.NewReceiptResolver(mockDecoder)

		_, err := resolver.Resolve(block)

		Expect(err).NotTo(HaveOccurred())
		mockDecoder.AssertDecodeCalledWith(test_helpers.FakeReceiptRLP, &types.Receipt{})
	})

	It("returns error if decoding fails", func() {
		block := blocks.NewBlock(test_helpers.FakeReceiptRLP)
		mockDecoder := rlp.NewMockDecoder()
		mockDecoder.SetError(test_helpers.FakeError)
		resolver := eth_block_receipts.NewReceiptResolver(mockDecoder)

		_, err := resolver.Resolve(block)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("returns receipt node for input", func() {
		block := blocks.NewBlock(test_helpers.FakeReceiptRLP)
		resolver := eth_block_receipts.NewReceiptResolver(rlp.NewMockDecoder())

		node, err := resolver.Resolve(block)

		Expect(err).NotTo(HaveOccurred())
		Expect(node).To(BeAssignableToTypeOf(&eth_block_receipts.EthReceiptNode{}))
		Expect(node.RawData()).To(Equal(test_helpers.FakeReceiptRLP))
		Expect(node.Cid()).To(Equal(block.Cid()))
	})
})
