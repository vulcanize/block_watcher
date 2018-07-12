package eth_block_header_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gx/ipfs/QmTRCUvZLiir12Qr6MV3HKfKMHX8Nf1Vddn6t2g5nsQSb9/go-block-format"
	"gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"

	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_header"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/go-ethereum/rlp"
)

var _ = Describe("Block header resolver", func() {
	It("decodes block's raw data into header", func() {
		cid := cid.NewCidV1(cid.EthBlock, test_helpers.FakeBlockRLP)
		block, err := blocks.NewBlockWithCid(test_helpers.FakeBlockRLP, cid)
		Expect(err).NotTo(HaveOccurred())
		mockDecoder := rlp.NewMockDecoder()
		resolver := eth_block_header.NewBlockHeaderResolver(mockDecoder)

		_, err = resolver.Resolve(block)

		Expect(err).NotTo(HaveOccurred())
		mockDecoder.AssertDecodeCalledWith(test_helpers.FakeBlockRLP, &types.Header{})
	})

	It("returns error if decoding fails", func() {
		cid := cid.NewCidV1(cid.EthBlock, test_helpers.FakeBlockRLP)
		block, err := blocks.NewBlockWithCid(test_helpers.FakeBlockRLP, cid)
		Expect(err).NotTo(HaveOccurred())
		mockDecoder := rlp.NewMockDecoder()
		mockDecoder.SetError(test_helpers.FakeError)
		resolver := eth_block_header.NewBlockHeaderResolver(mockDecoder)

		_, err = resolver.Resolve(block)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("returns eth block header node for input", func() {
		cid := cid.NewCidV1(cid.EthBlock, test_helpers.FakeBlockRLP)
		block, err := blocks.NewBlockWithCid(test_helpers.FakeBlockRLP, cid)
		Expect(err).NotTo(HaveOccurred())
		resolver := eth_block_header.NewBlockHeaderResolver(rlp.NewMockDecoder())

		node, err := resolver.Resolve(block)

		Expect(err).NotTo(HaveOccurred())
		Expect(node).To(BeAssignableToTypeOf(&eth_block_header.EthBlockHeaderNode{}))
		Expect(node.RawData()).To(Equal(test_helpers.FakeBlockRLP))
		Expect(node.Cid()).To(Equal(cid))
	})
})
