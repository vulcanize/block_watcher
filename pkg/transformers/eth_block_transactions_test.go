package transformers_test

import (
	"io/ioutil"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/eth-block-extractor/pkg/transformers"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/db"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
)

var _ = Describe("Eth block transactions transformer", func() {
	Describe("executing on a single block", func() {
		BeforeEach(func() {
			log.SetOutput(ioutil.Discard)
		})

		It("fetches rlp data for block body", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockBodyByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockPublisher := ipfs.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDB, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertGetBlockBodyByBlockNumberCalledWith([]int64{blockNumber})
		})

		It("returns error if fetching rlp returns error", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockBodyByBlockNumberError(test_helpers.FakeError)
			mockPublisher := ipfs.NewMockPublisher()
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDB, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(transformers.NewExecuteError(transformers.GetBlockRlpErr, test_helpers.FakeError)))
		})

		It("publishes block body data to IPFS", func() {
			mockDB := db.NewMockDatabase()
			fakeRawData := [][]byte{{1, 2, 3, 4, 5}}
			mockDB.SetGetBlockBodyByBlockNumberReturnBytes(fakeRawData)
			mockPublisher := ipfs.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDB, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockPublisher.AssertWriteCalledWithBytes(fakeRawData)
		})

		It("returns error if publishing data returns error", func() {
			mockDB := db.NewMockDatabase()
			fakeRawData := [][]byte{{1, 2, 3, 4, 5}}
			mockDB.SetGetBlockBodyByBlockNumberReturnBytes(fakeRawData)
			mockPublisher := ipfs.NewMockPublisher()
			mockPublisher.SetError(test_helpers.FakeError)
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDB, mockPublisher)
			blockNumber := int64(1234567)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(transformers.NewExecuteError(transformers.PutIpldErr, test_helpers.FakeError)))
		})
	})

	Describe("executing on a range of blocks", func() {
		It("fetches rlp data for every block's body", func() {
			mockDatabase := db.NewMockDatabase()
			mockDatabase.SetGetBlockBodyByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 0}})
			mockPublisher := ipfs.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid_one"}, {"cid_two"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDatabase, mockPublisher)
			startingBlockNumber := int64(1234567)
			endingBlockNumber := int64(1234568)

			err := transformer.Execute(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockDatabase.AssertGetBlockBodyByBlockNumberCalledWith([]int64{startingBlockNumber, endingBlockNumber})
		})

		It("publishes every block body's data to IPFS", func() {
			mockDatabase := db.NewMockDatabase()
			fakeRawData := [][]byte{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 0}}
			mockDatabase.SetGetBlockBodyByBlockNumberReturnBytes(fakeRawData)
			mockPublisher := ipfs.NewMockPublisher()
			mockPublisher.SetReturnStrings([][]string{{"cid_one"}, {"cid_two"}})
			transformer := transformers.NewEthBlockTransactionsTransformer(mockDatabase, mockPublisher)
			startingBlockNumber := int64(1234567)
			endingBlockNumber := int64(1234568)

			err := transformer.Execute(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockPublisher.AssertWriteCalledWithBytes(fakeRawData)
		})
	})
})
