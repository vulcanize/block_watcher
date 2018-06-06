package transformers_test

import (
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/block_watcher/pkg/transformers"
	"github.com/vulcanize/block_watcher/test_helpers"
)

var _ = Describe("Compute historical state transformer", func() {
	BeforeEach(func() {
		log.SetOutput(ioutil.Discard)
	})

	Describe("publishing the state trie for the genesis block", func() {
		It("fetches state trie root for genesis block", func() {
			db := test_helpers.NewMockDatabase()
			fakeHeaderBytes := []byte{1, 2, 3, 4, 5}
			db.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{fakeHeaderBytes})
			decoder := test_helpers.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{})
			publisher := test_helpers.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(db, decoder, publisher)

			err := transformer.Execute(0)

			Expect(err).NotTo(HaveOccurred())
			db.AssertGetBlockHeaderByBlockNumberCalledWith([]int64{0})
			decoder.AssertDecodeCalledWith(fakeHeaderBytes, &types.Header{})
		})

		It("returns error if fetching state trie root fails", func() {
			db := test_helpers.NewMockDatabase()
			db.SetGetBlockHeaderByBlockNumberError(test_helpers.FakeError)
			decoder := test_helpers.NewMockDecoder()
			publisher := test_helpers.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(db, decoder, publisher)

			err := transformer.Execute(0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
		})

		It("fetches state trie nodes for genesis block with block state root", func() {
			db := test_helpers.NewMockDatabase()
			db.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			decoder := test_helpers.NewMockDecoder()
			fakeRoot := common.HexToHash("0x123")
			decoder.SetReturnOut(&types.Header{Root: fakeRoot})
			publisher := test_helpers.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(db, decoder, publisher)

			err := transformer.Execute(0)

			Expect(err).NotTo(HaveOccurred())
			db.AssertGetStateTrieNodesCalledWith(fakeRoot.Bytes())
		})

		It("returns error if fetching state trie nodes fails", func() {
			db := test_helpers.NewMockDatabase()
			db.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			db.SetGetStateTrieNodesError(test_helpers.FakeError)
			decoder := test_helpers.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{})
			publisher := test_helpers.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(db, decoder, publisher)

			err := transformer.Execute(0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
		})

		It("publishes state trie nodes for genesis block to IPFS", func() {
			db := test_helpers.NewMockDatabase()
			db.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			fakeStateTrieNodes := [][]byte{{6, 7, 8, 9, 0}}
			db.SetGetStateTrieNodesReturnBytes(fakeStateTrieNodes)
			decoder := test_helpers.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: common.HexToHash("0x123")})
			publisher := test_helpers.NewMockPublisher()
			publisher.SetReturnStrings([][]string{{"one"}})
			transformer := transformers.NewComputeEthStateTrieTransformer(db, decoder, publisher)

			err := transformer.Execute(0)

			Expect(err).NotTo(HaveOccurred())
			publisher.AssertWriteCalledWith(fakeStateTrieNodes)
		})
	})

	Describe("computing and publishing the state trie for subsequent blocks", func() {
		It("fetches the current and parent block", func() {
			db := test_helpers.NewMockDatabase()
			db.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			db.SetGetStateTrieNodesReturnBytes([][]byte{{6, 7, 8, 9, 0}})
			decoder := test_helpers.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: common.HexToHash("0x123")})
			publisher := test_helpers.NewMockPublisher()
			publisher.SetReturnStrings([][]string{{"one"}})
			transformer := transformers.NewComputeEthStateTrieTransformer(db, decoder, publisher)

			err := transformer.Execute(4)

			Expect(err).NotTo(HaveOccurred())
			db.AssertGetBlockByBlockNumberCalledwith([]int64{0, 1, 2, 3, 4})
		})

		It("computes state trie nodes for current block", func() {
			db := test_helpers.NewMockDatabase()
			db.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			db.SetGetStateTrieNodesReturnBytes([][]byte{{6, 7, 8, 9, 0}})
			fakeBlock := &types.Block{}
			db.SetGetBlockByBlockNumberReturnBlock(fakeBlock)
			decoder := test_helpers.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: common.HexToHash("0x123")})
			publisher := test_helpers.NewMockPublisher()
			publisher.SetReturnStrings([][]string{{"one"}})
			transformer := transformers.NewComputeEthStateTrieTransformer(db, decoder, publisher)

			err := transformer.Execute(1)

			Expect(err).NotTo(HaveOccurred())
			db.AssertComputeBlockStateTrieCalledWith(fakeBlock, fakeBlock)
		})

		It("publishes state trie nodes to IPFS", func() {
			db := test_helpers.NewMockDatabase()
			db.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			genesisBlockStateTrieNodes := [][]byte{{6, 7, 8, 9, 0}}
			db.SetGetStateTrieNodesReturnBytes(genesisBlockStateTrieNodes)
			subsequentBlocksStateTrieNodes := [][]byte{{1, 1, 1, 1, 1}, {2, 2, 2, 2, 2}}
			db.SetComputeBlockStateTrieReturnBytes(subsequentBlocksStateTrieNodes)
			decoder := test_helpers.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: common.HexToHash("0x123")})
			publisher := test_helpers.NewMockPublisher()
			publisher.SetReturnStrings([][]string{{"one"}, {"two"}, {"three"}})
			transformer := transformers.NewComputeEthStateTrieTransformer(db, decoder, publisher)

			err := transformer.Execute(1)

			Expect(err).NotTo(HaveOccurred())
			publisher.AssertWriteCalledWith(append(genesisBlockStateTrieNodes, subsequentBlocksStateTrieNodes...))
		})
	})
})
