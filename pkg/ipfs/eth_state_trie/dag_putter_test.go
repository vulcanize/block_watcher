package eth_state_trie_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/block_watcher/pkg/ipfs/eth_state_trie"
	"github.com/vulcanize/block_watcher/test_helpers"
)

var _ = Describe("Ethereum state trie node dag putter", func() {
	It("adds passed state trie node to ipfs", func() {
		mockAdder := test_helpers.NewMockAdder()
		dagPutter := eth_state_trie.NewStateTrieDagPutter(mockAdder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).NotTo(HaveOccurred())
		mockAdder.AssertAddCalled(1, &eth_state_trie.EthStateTrieNode{})
	})

	It("returns error if adding to ipfs fails", func() {
		mockAdder := test_helpers.NewMockAdder()
		fakeError := errors.New("failed")
		mockAdder.SetError(fakeError)
		dagPutter := eth_state_trie.NewStateTrieDagPutter(mockAdder)

		_, err := dagPutter.DagPut([]byte{1, 2, 3, 4, 5})

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakeError))
	})
})
