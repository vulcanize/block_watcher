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
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/rlp"
)

var _ = Describe("Compute eth state trie transformer", func() {
	BeforeEach(func() {
		log.SetOutput(ioutil.Discard)
	})

	Describe("publishing the state trie for the genesis block", func() {
		It("fetches state trie root for genesis block", func() {
			mockDB := db.NewMockDatabase()
			fakeHeaderBytes := []byte{1, 2, 3, 4, 5}
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{fakeHeaderBytes})
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{})
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, ipfs.NewMockPublisher(), ipfs.NewMockPublisher())

			err := transformer.Execute(0)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertGetBlockHeaderByBlockNumberCalledWith([]int64{0})
			decoder.AssertDecodeCalledWith(fakeHeaderBytes, &types.Header{})
		})

		It("returns error if fetching state trie root fails", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberError(test_helpers.FakeError)
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, rlp.NewMockDecoder(), ipfs.NewMockPublisher(), ipfs.NewMockPublisher())

			err := transformer.Execute(0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
		})

		It("fetches state and storage trie nodes for genesis block with block state root", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: test_helpers.FakeHash})
			storageTriePublisher := ipfs.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, ipfs.NewMockPublisher(), storageTriePublisher)

			err := transformer.Execute(0)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertGetStateTrieNodesCalledWith(test_helpers.FakeHash)
		})

		It("returns error if fetching state trie nodes fails", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockDB.SetGetStateAndStorageTrieNodesError(test_helpers.FakeError)
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{})
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, ipfs.NewMockPublisher(), ipfs.NewMockPublisher())

			err := transformer.Execute(0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
		})

		It("publishes state trie nodes for genesis block to IPFS", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			fakeStateTrieNodes := [][]byte{{6, 7, 8, 9, 0}}
			mockDB.SetGetStateAndStorageTrieNodesReturnStateTrieBytes(fakeStateTrieNodes)
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: test_helpers.FakeHash})
			stateTriePublisher := ipfs.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, stateTriePublisher, ipfs.NewMockPublisher())

			err := transformer.Execute(0)

			Expect(err).NotTo(HaveOccurred())
			stateTriePublisher.AssertWriteCalledWithBytes(fakeStateTrieNodes)
		})

		It("returns error if publishing state trie nodes fails", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			fakeStateTrieNodes := [][]byte{{6, 7, 8, 9, 0}}
			mockDB.SetGetStateAndStorageTrieNodesReturnStateTrieBytes(fakeStateTrieNodes)
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: test_helpers.FakeHash})
			stateTriePublisher := ipfs.NewMockPublisher()
			stateTriePublisher.SetError(test_helpers.FakeError)
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, stateTriePublisher, ipfs.NewMockPublisher())

			err := transformer.Execute(0)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
		})
	})

	Describe("computing and publishing the state trie for subsequent blocks", func() {
		It("fetches the current and parent block", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockDB.SetGetStateAndStorageTrieNodesReturnStateTrieBytes([][]byte{{6, 7, 8, 9, 0}})
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: test_helpers.FakeHash})
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, ipfs.NewMockPublisher(), ipfs.NewMockPublisher())

			err := transformer.Execute(4)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertGetBlockByBlockNumberCalledwith([]int64{0, 1, 2, 3, 4})
		})

		It("computes state and storage trie nodes for current block", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockDB.SetGetStateAndStorageTrieNodesReturnStateTrieBytes([][]byte{{6, 7, 8, 9, 0}})
			fakeBlock := &types.Block{}
			mockDB.SetGetBlockByBlockNumberReturnBlock(fakeBlock)
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: test_helpers.FakeHash})
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, ipfs.NewMockPublisher(), ipfs.NewMockPublisher())

			err := transformer.Execute(1)

			Expect(err).NotTo(HaveOccurred())
			mockDB.AssertComputeBlockStateTrieCalledWith(fakeBlock, fakeBlock)
		})

		It("publishes state trie nodes to IPFS", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			fakeStateTrieNodes := [][]byte{{0, 0, 0, 0, 0}, {1, 1, 1, 1, 1}}
			mockDB.SetGetStateAndStorageTrieNodesReturnStateTrieBytes(fakeStateTrieNodes)
			mockDB.SetComputeBlockStateTrieReturnHash(test_helpers.FakeHash)
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: common.HexToHash("0x123")})
			stateTriePublisher := ipfs.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, stateTriePublisher, ipfs.NewMockPublisher())

			err := transformer.Execute(1)

			Expect(err).NotTo(HaveOccurred())
			stateTriePublisher.AssertWriteCalledWithBytes(fakeStateTrieNodes)
		})

		It("returns error if publishing state trie nodes fails", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockDB.SetGetStateAndStorageTrieNodesReturnStateTrieBytes([][]byte{{6, 7, 8, 9, 0}})
			mockDB.SetComputeBlockStateTrieReturnHash(test_helpers.FakeHash)
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: common.HexToHash("0x123")})
			stateTriePublisher := ipfs.NewMockPublisher()
			stateTriePublisher.SetError(test_helpers.FakeError)
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, stateTriePublisher, ipfs.NewMockPublisher())

			err := transformer.Execute(1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
		})

		It("publishes storage trie nodes to IPFS", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockDB.SetGetStateAndStorageTrieNodesReturnStateTrieBytes(test_helpers.FakeTrieNodes)
			fakeStorageTrieNodes := [][]byte{{2, 2, 2, 2, 2}}
			mockDB.SetGetStateAndStorageTrieNodesReturnStorageTrieBytes(fakeStorageTrieNodes)
			mockDB.SetComputeBlockStateTrieReturnHash(test_helpers.FakeHash)
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: common.HexToHash("0x123")})
			storageTriePublisher := ipfs.NewMockPublisher()
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, ipfs.NewMockPublisher(), storageTriePublisher)

			err := transformer.Execute(1)

			Expect(err).NotTo(HaveOccurred())
			storageTriePublisher.AssertWriteCalledWithBytes(fakeStorageTrieNodes)
		})

		It("returns error if publishing storage trie nodes fails", func() {
			mockDB := db.NewMockDatabase()
			mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
			mockDB.SetGetStateAndStorageTrieNodesReturnStateTrieBytes(test_helpers.FakeTrieNodes)
			mockDB.SetGetStateAndStorageTrieNodesReturnStorageTrieBytes(test_helpers.FakeTrieNodes)
			mockDB.SetComputeBlockStateTrieReturnHash(test_helpers.FakeHash)
			decoder := rlp.NewMockDecoder()
			decoder.SetReturnOut(&types.Header{Root: common.HexToHash("0x123")})
			storageTriePublisher := ipfs.NewMockPublisher()
			storageTriePublisher.SetError(test_helpers.FakeError)
			transformer := transformers.NewComputeEthStateTrieTransformer(mockDB, decoder, ipfs.NewMockPublisher(), storageTriePublisher)

			err := transformer.Execute(1)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
		})
	})
})
