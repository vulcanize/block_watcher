package factories_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_header"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_receipts"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/eth_block_transactions"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/factories"
	"github.com/vulcanize/eth-block-extractor/pkg/ipfs/util"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"gx/ipfs/QmapdYm1b22Frv3k17fqrBYTFRxwiaVJkB299Mfn33edeB/go-cid"
)

var _ = Describe("Resolver factory", func() {
	It("returns block header resolver for block header cid", func() {
		blockHeaderCid, err := util.RawToCid(cid.EthBlock, test_helpers.FakeBlockRLP)
		Expect(err).NotTo(HaveOccurred())
		factory := factories.ResolverFactory{}

		resolver, err := factory.GetResolver(blockHeaderCid)

		Expect(err).NotTo(HaveOccurred())
		Expect(resolver).To(BeAssignableToTypeOf(eth_block_header.BlockHeaderResolver{}))
	})

	It("returns transaction resolver for transaction cid", func() {
		transactionCid, err := util.RawToCid(cid.EthTx, test_helpers.FakeTransactionRLP)
		Expect(err).NotTo(HaveOccurred())
		factory := factories.ResolverFactory{}

		resolver, err := factory.GetResolver(transactionCid)

		Expect(err).NotTo(HaveOccurred())
		Expect(resolver).To(BeAssignableToTypeOf(eth_block_transactions.TransactionResolver{}))
	})

	It("returns receipt resolver for receipt cid", func() {
		receiptCid, err := util.RawToCid(cid.EthTxReceipt, test_helpers.FakeReceiptRLP)
		Expect(err).NotTo(HaveOccurred())
		factory := factories.ResolverFactory{}

		resolver, err := factory.GetResolver(receiptCid)

		Expect(err).NotTo(HaveOccurred())
		Expect(resolver).To(BeAssignableToTypeOf(eth_block_receipts.ReceiptResolver{}))
	})

	It("returns error for unknown cid", func() {
		unknownCid, err := util.RawToCid(cid.BitcoinTx, test_helpers.FakeTransactionRLP)
		Expect(err).NotTo(HaveOccurred())
		factory := factories.ResolverFactory{}

		_, err = factory.GetResolver(unknownCid)

		Expect(err).To(HaveOccurred())
	})
})
