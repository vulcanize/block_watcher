package eth_block_transactions_test

import (
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/block_watcher/pkg/ipfs/eth_block_transactions"
	"github.com/vulcanize/block_watcher/test_helpers"
	"github.com/vulcanize/block_watcher/test_helpers/mocks/db"
	"github.com/vulcanize/block_watcher/test_helpers/mocks/ipfs"
)

var _ = Describe("Eth block transactions dag putter", func() {
	It("decodes passed raw data into an ethereum block body", func() {
		mockAdder := ipfs.NewMockAdder()
		mockDecoder := db.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Body{})
		dagPutter := eth_block_transactions.NewBlockTransactionsDagPutter(mockAdder, mockDecoder)
		fakeBytes := []byte{1, 2, 3, 4, 5}

		_, err := dagPutter.DagPut(fakeBytes)

		Expect(err).NotTo(HaveOccurred())
		mockDecoder.AssertDecodeCalledWith(fakeBytes, &types.Body{})
	})

	It("returns error if decoding fails", func() {
		mockAdder := ipfs.NewMockAdder()
		mockDecoder := db.NewMockDecoder()
		mockDecoder.SetError(test_helpers.FakeError)
		dagPutter := eth_block_transactions.NewBlockTransactionsDagPutter(mockAdder, mockDecoder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("adds a node for each transaction on the block", func() {
		mockAdder := ipfs.NewMockAdder()
		mockDecoder := db.NewMockDecoder()
		fakeTransactionOne := &types.Transaction{}
		fakeTransactionTwo := &types.Transaction{}
		fakeBlockBody := &types.Body{
			Transactions: types.Transactions{fakeTransactionOne, fakeTransactionTwo},
			Uncles:       nil,
		}
		mockDecoder.SetReturnOut(fakeBlockBody)
		dagPutter := eth_block_transactions.NewBlockTransactionsDagPutter(mockAdder, mockDecoder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).NotTo(HaveOccurred())
		mockAdder.AssertAddCalled(2, &eth_block_transactions.EthTransactionNode{})
	})

	It("returns error if adding node fails", func() {
		mockAdder := ipfs.NewMockAdder()
		mockAdder.SetError(test_helpers.FakeError)
		mockDecoder := db.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Body{Transactions: types.Transactions{&types.Transaction{}}})
		dagPutter := eth_block_transactions.NewBlockTransactionsDagPutter(mockAdder, mockDecoder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})
})
