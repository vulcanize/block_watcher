package level_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/db/level"
	"github.com/vulcanize/eth-block-extractor/test_helpers"
	level_wrapper "github.com/vulcanize/eth-block-extractor/test_helpers/mocks/db/level"
	"github.com/vulcanize/eth-block-extractor/test_helpers/mocks/wrappers/core/rawdb"
)

var _ = Describe("Database", func() {
	Describe("Computing state trie nodes", func() {
		It("invokes state computer to build historical state", func() {
			mockStateComputer := level_wrapper.NewMockStateComputer()
			mockStateTrieReader := level_wrapper.NewMockStateTrieReader()
			db := level.NewLevelDatabase(rawdb.NewMockAccessorsChain(), mockStateComputer, mockStateTrieReader)
			currentBlock := &types.Block{}
			parentBlock := &types.Block{}

			_, err := db.ComputeBlockStateTrie(currentBlock, parentBlock)

			Expect(err).NotTo(HaveOccurred())
			mockStateComputer.AssertComputeBlockStateTrieCalledWith(currentBlock, parentBlock)
		})

		It("returns err if state computer returns err", func() {
			mockStateComputer := level_wrapper.NewMockStateComputer()
			mockStateComputer.SetComputeBlockStateTrieReturnErr(test_helpers.FakeError)
			mockStateTrieReader := level_wrapper.NewMockStateTrieReader()
			db := level.NewLevelDatabase(rawdb.NewMockAccessorsChain(), mockStateComputer, mockStateTrieReader)

			_, err := db.ComputeBlockStateTrie(&types.Block{}, &types.Block{})

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(test_helpers.FakeError))
		})
	})

	Describe("Getting block body data", func() {
		It("invokes the level database reader to query for block hash by block number", func() {
			mockLevelDBReader := rawdb.NewMockAccessorsChain()
			mockStateTrieReader := level_wrapper.NewMockStateTrieReader()
			db := level.NewLevelDatabase(mockLevelDBReader, level_wrapper.NewMockStateComputer(), mockStateTrieReader)
			num := int64(123456)

			_, err := db.GetBlockBodyByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the level database reader to query for block body data", func() {
			mockLevelDBReader := rawdb.NewMockAccessorsChain()
			hash := common.HexToHash("abcde")
			mockLevelDBReader.SetGetCanonicalHashReturnHash(hash)
			mockStateTrieReader := level_wrapper.NewMockStateTrieReader()
			db := level.NewLevelDatabase(mockLevelDBReader, level_wrapper.NewMockStateComputer(), mockStateTrieReader)
			num := int64(123456)

			_, err := db.GetBlockBodyByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetBodyRLPCalledWith(hash, uint64(num))
		})
	})

	Describe("Getting block", func() {
		It("invokes the level database reader to query for block hash by block number", func() {
			mockLevelDBReader := rawdb.NewMockAccessorsChain()
			mockStateTrieReader := level_wrapper.NewMockStateTrieReader()
			db := level.NewLevelDatabase(mockLevelDBReader, level_wrapper.NewMockStateComputer(), mockStateTrieReader)
			num := int64(123456)

			db.GetBlockByBlockNumber(num)

			mockLevelDBReader.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the level database reader to query for block", func() {
			mockLevelDBReader := rawdb.NewMockAccessorsChain()
			hash := common.HexToHash("abcde")
			mockLevelDBReader.SetGetCanonicalHashReturnHash(hash)
			mockStateTrieReader := level_wrapper.NewMockStateTrieReader()
			db := level.NewLevelDatabase(mockLevelDBReader, level_wrapper.NewMockStateComputer(), mockStateTrieReader)
			num := int64(123456)

			db.GetBlockByBlockNumber(num)

			mockLevelDBReader.AssertGetBlockCalledWith(hash, uint64(num))
		})
	})

	Describe("Getting block header data", func() {
		It("invokes the level database reader to query for block hash by block number", func() {
			mockLevelDBReader := rawdb.NewMockAccessorsChain()
			mockStateTrieReader := level_wrapper.NewMockStateTrieReader()
			db := level.NewLevelDatabase(mockLevelDBReader, level_wrapper.NewMockStateComputer(), mockStateTrieReader)
			num := int64(123456)

			_, err := db.GetBlockHeaderByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the level database reader to query for block header data", func() {
			mockLevelDBReader := rawdb.NewMockAccessorsChain()
			hash := common.HexToHash("abcde")
			mockLevelDBReader.SetGetCanonicalHashReturnHash(hash)
			mockStateTrieReader := level_wrapper.NewMockStateTrieReader()
			db := level.NewLevelDatabase(mockLevelDBReader, level_wrapper.NewMockStateComputer(), mockStateTrieReader)
			num := int64(123456)

			_, err := db.GetBlockHeaderByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetHeaderRLPCalledWith(hash, uint64(num))
		})
	})

	Describe("Getting state trie nodes", func() {
		It("invokes the level database reader to query for state trie data", func() {
			mockLevelDBReader := rawdb.NewMockAccessorsChain()
			mockStateTrieReader := level_wrapper.NewMockStateTrieReader()
			db := level.NewLevelDatabase(mockLevelDBReader, level_wrapper.NewMockStateComputer(), mockStateTrieReader)
			root := common.HexToHash("abcde")

			_, _, err := db.GetStateAndStorageTrieNodes(root)

			Expect(err).NotTo(HaveOccurred())
			mockStateTrieReader.AssertGetStateAndStorageTrieNodesCalledWith(root)
		})
	})
})
