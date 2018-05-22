package level_test

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/block_watcher/pkg/db/level"
	"github.com/vulcanize/block_watcher/test_helpers"
)

var _ = Describe("Database", func() {
	Describe("Getting block body data", func() {
		It("invokes the level database reader to query for block hash by block number", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			db := level.NewLevelDatabase(mockLevelDBReader)
			num := int64(123456)

			_, err := db.GetBlockBodyByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the level database reader to query for block body data", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			hash := common.HexToHash("abcde")
			mockLevelDBReader.SetGetCanonicalHashReturnHash(hash)
			db := level.NewLevelDatabase(mockLevelDBReader)
			num := int64(123456)

			_, err := db.GetBlockBodyByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetBodyRLPCalledWith(hash, uint64(num))
		})
	})

	Describe("Getting block header data", func() {
		It("invokes the level database reader to query for block hash by block number", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			db := level.NewLevelDatabase(mockLevelDBReader)
			num := int64(123456)

			_, err := db.GetBlockHeaderByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetCanonicalHashCalledWith(uint64(num))
		})

		It("invokes the level database reader to query for block header data", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			hash := common.HexToHash("abcde")
			mockLevelDBReader.SetGetCanonicalHashReturnHash(hash)
			db := level.NewLevelDatabase(mockLevelDBReader)
			num := int64(123456)

			_, err := db.GetBlockHeaderByBlockNumber(num)

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetHeaderRLPCalledWith(hash, uint64(num))
		})
	})

	Describe("Getting state trie nodes", func() {
		It("invokes the level database reader to query for state trie data", func() {
			mockLevelDBReader := test_helpers.NewMockLevelDatabaseReader()
			db := level.NewLevelDatabase(mockLevelDBReader)
			root := common.HexToHash("abcde")

			_, err := db.GetStateTrieNodes(root.Bytes())

			Expect(err).NotTo(HaveOccurred())
			mockLevelDBReader.AssertGetStateTrieNodesCalledWith(root)
		})
	})
})
