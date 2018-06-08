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

var _ = Describe("Ethereum state trie transformer", func() {
	BeforeEach(func() {
		log.SetOutput(ioutil.Discard)
	})

	It("fetches block header for block", func() {
		mockDB := db.NewMockDatabase()
		mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
		mockDecoder := db.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		transformer := transformers.NewEthStateTrieTransformer(mockDB, mockDecoder, ipfs.NewMockPublisher())

		err := transformer.Execute(0, 0)

		Expect(err).NotTo(HaveOccurred())
		mockDB.AssertGetBlockHeaderByBlockNumberCalledWith([]int64{0})
	})

	It("returns err if fetching block header returns err", func() {
		mockDB := db.NewMockDatabase()
		mockDB.SetGetBlockHeaderByBlockNumberError(test_helpers.FakeError)
		transformer := transformers.NewEthStateTrieTransformer(mockDB, db.NewMockDecoder(), ipfs.NewMockPublisher())

		err := transformer.Execute(0, 0)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
	})

	It("fetches state trie nodes with state root from decoded block header", func() {
		mockDB := db.NewMockDatabase()
		fakeHeader := []byte{6, 7, 8, 9, 0}
		mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{fakeHeader})
		mockDecoder := db.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{Root: test_helpers.FakeHash})
		transformer := transformers.NewEthStateTrieTransformer(mockDB, mockDecoder, ipfs.NewMockPublisher())

		err := transformer.Execute(0, 0)

		Expect(err).NotTo(HaveOccurred())
		mockDecoder.AssertDecodeCalledWith(fakeHeader, &types.Header{})
		mockDB.AssertGetStateTrieNodesCalledWith(test_helpers.FakeHash.Bytes())
	})

	It("returns err if fetching state trie returns err", func() {
		mockDB := db.NewMockDatabase()
		mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
		mockDB.SetGetStateTrieNodesError(test_helpers.FakeError)
		mockDecoder := db.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{Root: common.Hash{}})
		transformer := transformers.NewEthStateTrieTransformer(mockDB, mockDecoder, ipfs.NewMockPublisher())

		err := transformer.Execute(0, 0)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(test_helpers.FakeError.Error()))
	})

	It("writes state trie nodes to ipfs", func() {
		mockDB := db.NewMockDatabase()
		mockDB.SetGetBlockHeaderByBlockNumberReturnBytes([][]byte{{1, 2, 3, 4, 5}})
		fakeStateTrieNodes := [][]byte{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 0}}
		mockDB.SetGetStateTrieNodesReturnBytes(fakeStateTrieNodes)
		mockDecoder := db.NewMockDecoder()
		mockDecoder.SetReturnOut(&types.Header{})
		mockPublisher := ipfs.NewMockPublisher()
		mockPublisher.SetReturnStrings([][]string{{"one"}, {"two"}})
		transformer := transformers.NewEthStateTrieTransformer(mockDB, mockDecoder, mockPublisher)

		err := transformer.Execute(0, 0)

		Expect(err).NotTo(HaveOccurred())
		mockPublisher.AssertWriteCalledWith(fakeStateTrieNodes)
	})
})
