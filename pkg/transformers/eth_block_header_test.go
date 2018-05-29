package transformers_test

import (
	"errors"
	"io/ioutil"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/block_watcher/pkg/transformers"
	"github.com/vulcanize/block_watcher/test_helpers"
)

var _ = Describe("EthBlockHeaderTransformer", func() {
	var mockDatabase *test_helpers.MockDatabase
	var mockPublisher *test_helpers.MockPublisher

	Describe("Fetching one block header", func() {
		var fakeBytes []byte
		var blockNumber int64

		BeforeEach(func() {
			mockDatabase = test_helpers.NewMockDatabase()
			mockPublisher = test_helpers.NewMockPublisher()
			blockNumber = 54321
			fakeBytes = []byte{6, 7, 8, 9, 0}
			mockDatabase.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{fakeBytes})
			mockPublisher.SetReturnStrings([][]string{{"cid_one", "cid_two"}})
			log.SetOutput(ioutil.Discard)
		})

		It("Fetches RLP data from ethereum db", func() {
			transformer := transformers.NewEthBlockHeaderTransformer(mockDatabase, mockPublisher)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockDatabase.AssertGetBlockHeaderByBlockNumberCalledWith([]int64{blockNumber})
		})

		It("Returns error if fetching block RLP data from ethereum DB fails", func() {
			fakeError := errors.New("failed")
			mockDatabase.SetGetBlockHeaderByBlockNumberError(fakeError)
			transformer := transformers.NewEthBlockHeaderTransformer(mockDatabase, mockPublisher)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(transformers.NewExecuteError(transformers.GetBlockRlpErr, fakeError)))
			mockDatabase.AssertGetBlockHeaderByBlockNumberCalledWith([]int64{blockNumber})
		})

		It("Persists block RLP data to IPFS", func() {
			transformer := transformers.NewEthBlockHeaderTransformer(mockDatabase, mockPublisher)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockPublisher.AssertWriteCalledWith([][]byte{fakeBytes})
		})

		It("Returns err if persisting block RLP data to IPFS fails", func() {
			fakeError := errors.New("failed")
			mockPublisher.SetError(fakeError)
			transformer := transformers.NewEthBlockHeaderTransformer(mockDatabase, mockPublisher)

			err := transformer.Execute(blockNumber, blockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(transformers.NewExecuteError(transformers.PutIpldErr, fakeError)))
			mockDatabase.AssertGetBlockHeaderByBlockNumberCalledWith([]int64{blockNumber})
			mockPublisher.AssertWriteCalledWith([][]byte{fakeBytes})
		})
	})

	Describe("Fetching multiple blocks", func() {
		var fakeRlpBytes []byte
		var startingBlockNumber int64
		var endingBlockNumber int64

		BeforeEach(func() {
			mockDatabase = test_helpers.NewMockDatabase()
			mockPublisher = test_helpers.NewMockPublisher()
			startingBlockNumber = 54321
			endingBlockNumber = 54322
			fakeRlpBytes = []byte{6, 7, 8, 9, 0}
			mockDatabase.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{fakeRlpBytes, fakeRlpBytes})
			fakeOutputString := []string{"cid_one", "cid_two"}
			mockPublisher.SetReturnStrings([][]string{fakeOutputString, fakeOutputString})
			log.SetOutput(ioutil.Discard)
		})

		It("Fetches block RLP data from ethereum db for every block in range", func() {
			transformer := transformers.NewEthBlockHeaderTransformer(mockDatabase, mockPublisher)

			err := transformer.Execute(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockDatabase.AssertGetBlockHeaderByBlockNumberCalledWith([]int64{startingBlockNumber, endingBlockNumber})
		})

		It("Persists block RLP data to IPFS for every block in range", func() {
			transformer := transformers.NewEthBlockHeaderTransformer(mockDatabase, mockPublisher)

			err := transformer.Execute(startingBlockNumber, endingBlockNumber)

			Expect(err).NotTo(HaveOccurred())
			mockPublisher.AssertWriteCalledWith([][]byte{fakeRlpBytes, fakeRlpBytes})
		})
	})

})
