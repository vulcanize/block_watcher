package transformers_test

import (
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/eth-block-extractor/pkg/transformers"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/db"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/ipfs"
)

var _ = Describe("Compute historical state transformer", func() {
	BeforeEach(func() {
		log.SetOutput(ioutil.Discard)
	})

	Describe("publishing the state trie for the genesis block", func() {
		It("fetches state trie root for genesis block", func() {
			mockDB := db.NewMockDatabase()
			fakeHeaderBytes := []byte{1, 2, 3, 4, 5}
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{fakeHeaderBytes})
			decoder := db.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{})
			publisher := ipfs.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, publisher)

			err := transformer.Execute(0)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertGetBlockHeaderByBlockNumberCalledWith([]int64{0})
			decoder.AssertDecodeCalledWith(fakeHeaderBytes, &types.Header{})
		})

		It("returns error if fetching state trie root fails", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberError(test_helpers.FakeError)
			decoder := db.NewMockDecoder()
			publisher := ipfs.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, publisher)

			err := transformer.Execute(0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
		})

		It("fetches state trie nodes for genesis block with block state root", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			decoder := db.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: test_helpers.FakeHash})
			publisher := ipfs.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, publisher)

			err := transformer.Execute(0)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertGetStateTrieNodesCalledWith(test_helpers.FakeHash.Bytes())
		})

		It("returns error if fetching state trie nodes fails", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockDB.SetGetStateTrieNodesError(test_helpers.FakeError)
			decoder := db.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{})
			publisher := ipfs.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, publisher)

			err := transformer.Execute(0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
		})

		It("publishes state trie nodes for genesis block to IPFS", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			fakeStateTrieNodes := [][]byte{{6, 7, 8, 9, 0}}
			mockDB.SetGetStateTrieNodesReturnBytes(fakeStateTrieNodes)
			decoder := db.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: test_helpers.FakeHash})
			publisher := ipfs.NewMockPublisher()
			publisher.SetReturnStrings([][]string{{"one"}})
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, publisher)

			err := transformer.Execute(0)

			Expect(err).NotTo(HaveOccurred())
			publisher.AssertWriteCalledWith(fakeStateTrieNodes)
		})
	})

	Describe("computing and publishing the state trie for subsequent blocks", func() {
		It("fetches the current and parent block", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockDB.SetGetStateTrieNodesReturnBytes([][]byte{{6, 7, 8, 9, 0}})
			decoder := db.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: test_helpers.FakeHash})
			publisher := ipfs.NewMockPublisher()
			publisher.SetReturnStrings([][]string{{"one"}})
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, publisher)

			err := transformer.Execute(4)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertGetBlockByBlockNumberCalledwith([]int64{0, 1, 2, 3, 4})
		})

		It("computes state trie nodes for current block", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockDB.SetGetStateTrieNodesReturnBytes([][]byte{{6, 7, 8, 9, 0}})
			fakeBlock := &types.Block{}
			mockDB.SetGetBlockByBlockNumberReturnBlock(fakeBlock)
			decoder := db.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: test_helpers.FakeHash})
			publisher := ipfs.NewMockPublisher()
			publisher.SetReturnStrings([][]string{{"one"}})
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, publisher)

			err := transformer.Execute(1)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertComputeBlockStateTrieCalledWith(fakeBlock, fakeBlock)
		})

		It("publishes state trie nodes to IPFS", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			genesisBlockStateTrieNodes := [][]byte{{6, 7, 8, 9, 0}}
			mockDB.SetGetStateTrieNodesReturnBytes(genesisBlockStateTrieNodes)
			subsequentBlocksStateTrieNodes := [][]byte{{1, 1, 1, 1, 1}, {2, 2, 2, 2, 2}}
			mockDB.SetComputeBlockStateTrieReturnBytes(subsequentBlocksStateTrieNodes)
			decoder := db.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: common.HexToHash("0x123")})
			publisher := ipfs.NewMockPublisher()
			publisher.SetReturnStrings([][]string{{"one"}, {"two"}, {"three"}})
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, publisher)

			err := transformer.Execute(1)

			Expect(err).NotTo(HaveOccurred())
			publisher.AssertWriteCalledWith(append(genesisBlockStateTrieNodes, subsequentBlocksStateTrieNodes...))
		})
	})
})
