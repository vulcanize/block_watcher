package transformers_test

import (
	"io/ioutil"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/block_watcher/pkg/transformers"
	"github.com/vulcanize/block_watcher/test_helpers"
)

var _ = Describe("Eth block transactions transformer", func() {
	Describe("executing on a single block", func() {
		BeforeEach(func() {
			log.SetOutput(ioutil.Discard)
		})

		It("fetches rlp data for block body", func() {
			mockDatabase := test_helpers.NewMockDatabase()
			mockDatabase.SetGetBlockBodyByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockPublisher := test_helpers.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDatabase, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockDatabase.AssertGetBlockBodyByBlockNumberCalledWith([]int64{blockNumber})
		})

		It("returns error if fetching rlp returns error", func() {
			mockDatabase := test_helpers.NewMockDatabase()
			mockDatabase.SetGetBlockBodyByBlockNumberError(test_helpers.FakeError)
			mockPublisher := test_helpers.NewMockPublisher()
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDatabase, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(transformers.NewExecuteError(transformers.GetBlockRlpErr, test_helpers.FakeError)))
		})

		It("publishes block body data to IPFS", func() {
			mockDatabase := test_helpers.NewMockDatabase()
			fakeRawData := [][]byte{{1, 2, 3, 4, 5}}
			mockDatabase.SetGetBlockBodyByBlockNumberReturnBytes(fakeRawData)
			mockPublisher := test_helpers.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDatabase, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockPublisher.AssertWriteCalledWith(fakeRawData)
		})

		It("returns error if publishing data returns error", func() {
			mockDatabase := test_helpers.NewMockDatabase()
			fakeRawData := [][]byte{{1, 2, 3, 4, 5}}
			mockDatabase.SetGetBlockBodyByBlockNumberReturnBytes(fakeRawData)
			mockPublisher := test_helpers.NewMockPublisher()
			mockPublisher.SetError(test_helpers.FakeError)
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDatabase, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(transformers.NewExecuteError(transformers.PutIpldErr, test_helpers.FakeError)))
		})
	})

	Describe("executing on a range of blocks", func() {
		It("fetches rlp data for every block's body", func() {
			mockDatabase := test_helpers.NewMockDatabase()
			mockDatabase.SetGetBlockBodyByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 0}})
			mockPublisher := test_helpers.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid_one"}, {"cid_two"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDatabase, mockPublisher)
			startingBlockNumber := int64(1234567)
			endingBlockNumber := int64(1234568)

			err := transformer.Execute(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockDatabase.AssertGetBlockBodyByBlockNumberCalledWith([]int64{startingBlockNumber, endingBlockNumber})
		})

		It("publishes every block body's data to IPFS", func() {
			mockDatabase := test_helpers.NewMockDatabase()
			fakeRawData := [][]byte{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 0}}
			mockDatabase.SetGetBlockBodyByBlockNumberReturnBytes(fakeRawData)
			mockPublisher := test_helpers.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid_one"}, {"cid_two"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDatabase, mockPublisher)
			startingBlockNumber := int64(1234567)
			endingBlockNumber := int64(1234568)

			err := transformer.Execute(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockPublisher.AssertWriteCalledWith(fakeRawData)
		})
	})
})
