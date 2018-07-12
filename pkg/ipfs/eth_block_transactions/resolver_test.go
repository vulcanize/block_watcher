package eth_block_transactions_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_transactions"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/go-ethereum/rlp"
	"gx/ipfs/QmTRCUvZLiir12Qr6MV3HKfKMHX8Nf1Vddn6t2g5nsQSb9/go-block-format"
)

var _ = Describe("Eth block transactions resolver", func() {
	It("decodes block's raw data into transaction", func() {
		decoder := rlp.NewMockDecoder()
		resolver := eth_block_transactions.NewTransactionResolver(decoder)
		transactionBlock := blocks.NewBlock(test_helpers.FakeTransactionRLP)

		_, err := resolver.Resolve(transactionBlock)

		Expect(err).NotTo(HaveOccurred())
		decoder.AssertDecodeCalledWith(test_helpers.FakeTransactionRLP, &types.Transaction{})
	})

	It("returns error if decoding fails", func() {
		decoder := rlp.NewMockDecoder()
		decoder.SetError(test_helpers.FakeError)
		resolver := eth_block_transactions.NewTransactionResolver(decoder)
		transactionBlock := blocks.NewBlock(test_helpers.FakeTransactionRLP)

		_, err := resolver.Resolve(transactionBlock)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("returns eth transaction node for input", func() {
		decoder := rlp.NewMockDecoder()
		resolver := eth_block_transactions.NewTransactionResolver(decoder)
		transactionBlock := blocks.NewBlock(test_helpers.FakeTransactionRLP)

		node, err := resolver.Resolve(transactionBlock)

		Expect(err).NotTo(HaveOccurred())
		Expect(node).To(BeAssignableToTypeOf(&eth_block_transactions.EthTransactionNode{}))
		Expect(node.RawData()).To(Equal(test_helpers.FakeTransactionRLP))
		Expect(node.Cid()).To(Equal(transactionBlock.Cid()))
	})
})
