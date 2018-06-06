package level_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/block_watcher/pkg/db/level"
	"github.com/vulcanize/block_watcher/test_helpers"
)

var _ = Describe("Database", func() {
	Describe("Computing state trie nodes", func() {
		It("invokes state computer to build historical state", func() {
			mockStateComputer := test_helpers.NewMockStateComputer()
			db := level.NewLevelDatabase(test_helpers.NewMockLevelDatabaseReader(), mockStateComputer)
			currentBlock := &types.Block{}
			parentBlock := &types.Block{}

			_, err := db.ComputeBlockStateTrie(currentBlock, parentBlock)

			Expect(err).NotTo(HaveOccurred())
			mockStateComputer.AssertComputeBlockStateTrieCalledWith(currentBlock, parentBlock)
		})

		It("returns err if state computer returns err", func() {
			mockStateComputer := test_helpers.NewMockStateComputer()
			mockStateComputer.SetComputeBlockStateTrieReturnErr(test_helpers.FakeError)
			db := level.NewLevelDatabase(test_helpers.NewMockLevelDatabaseReader(), mockStateComputer)

			_, err := db.ComputeBlockStateTrie(&types.Block{}, &types.Block{})

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(test_helpers.FakeError))
		})
	})

	Describe("Getting block body data", func() {
		It("invokes the level database reader to query for block hash by block number", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			db := level.NewLevelDatabase(mockLevelDBReader, test_helpers.NewMockStateComputer())
			num := int64(123456)

			_, err := db.GetBlockBodyByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the level database reader to query for block body data", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			hash := common.HexToHash("abcde")
			mockLevelDBReader.SetGetCanonicalHashReturnHash(hash)
			db := level.NewLevelDatabase(mockLevelDBReader, test_helpers.NewMockStateComputer())
			num := int64(123456)

			_, err := db.GetBlockBodyByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetBodyRLPCalledWith(hash, uint64(num))
		})
	})

	Describe("Getting block", func() {
		It("invokes the level database reader to query for block hash by block number", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			db := level.NewLevelDatabase(mockLevelDBReader, test_helpers.NewMockStateComputer())
			num := int64(123456)

			db.GetBlockByBlockNumber(num)

			mockLevelDBReader.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the level database reader to query for block", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			hash := common.HexToHash("abcde")
			mockLevelDBReader.SetGetCanonicalHashReturnHash(hash)
			db := level.NewLevelDatabase(mockLevelDBReader, test_helpers.NewMockStateComputer())
			num := int64(123456)

			db.GetBlockByBlockNumber(num)

			mockLevelDBReader.AssertGetBlockCalledWith(hash, uint64(num))
		})
	})

	Describe("Getting block header data", func() {
		It("invokes the level database reader to query for block hash by block number", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			db := level.NewLevelDatabase(mockLevelDBReader, test_helpers.NewMockStateComputer())
			num := int64(123456)

			_, err := db.GetBlockHeaderByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the level database reader to query for block header data", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			hash := common.HexToHash("abcde")
			mockLevelDBReader.SetGetCanonicalHashReturnHash(hash)
			db := level.NewLevelDatabase(mockLevelDBReader, test_helpers.NewMockStateComputer())
			num := int64(123456)

			_, err := db.GetBlockHeaderByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetHeaderRLPCalledWith(hash, uint64(num))
		})
	})

	Describe("Getting state trie nodes", func() {
		It("invokes the level database reader to query for state trie data", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			db := level.NewLevelDatabase(mockLevelDBReader, test_helpers.NewMockStateComputer())
			root := common.HexToHash("abcde")

			_, err := db.GetStateTrieNodes(root.Bytes())

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetStateTrieNodesCalledWith(root)
		})
	})
})
