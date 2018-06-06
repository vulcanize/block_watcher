package state_computation_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/block_watcher/pkg/db/level/state_computation"
	"github.com/vulcanize/block_watcher/test_helpers"
	"github.com/vulcanize/block_watcher/test_helpers/state_mocks"
)

var _ = Describe("", func() {
	It("initializes state trie at parent block's root", func() {
		chain, db, iteratorFactory, processor, trieFactory, validator := getMocks()
		fakeDB := db.CreateFakeUnderlyingDatabase()
		db.SetReturnDatabase(fakeDB)
		computer := state_computation.NewStateComputer(chain, db, iteratorFactory, processor, trieFactory, validator)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).NotTo(HaveOccurred())
		trieFactory.AssertNewStateTrieCalledWith(parentBlock.Root(), fakeDB)
	})

	It("returns error if state trie initialization fails", func() {
		chain, db, iteratorFactory, processor, trieFactory, validator := getMocks()
		trieFactory.SetReturnErr(test_helpers.FakeError)
		computer := state_computation.NewStateComputer(chain, db, iteratorFactory, processor, trieFactory, validator)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("processes the block to build the state trie", func() {
		chain, db, iteratorFactory, processor, trieFactory, validator := getMocks()
		computer := state_computation.NewStateComputer(chain, db, iteratorFactory, processor, trieFactory, validator)
		stateTrie := state_mocks.NewMockStateTrie()
		fakeStateDB := &state.StateDB{}
		stateTrie.SetStateDb(fakeStateDB)
		trieFactory.SetStateTrie(stateTrie)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).NotTo(HaveOccurred())
		processor.AssertProcessCalledWith(currentBlock, fakeStateDB)
	})

	It("returns error if processing block fails", func() {
		chain, db, iteratorFactory, processor, trieFactory, validator := getMocks()
		processor.SetReturnErr(test_helpers.FakeError)
		computer := state_computation.NewStateComputer(chain, db, iteratorFactory, processor, trieFactory, validator)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("validates state computed by processing blocks", func() {
		chain, db, iteratorFactory, processor, trieFactory, validator := getMocks()
		fakeReceipts := types.Receipts{}
		processor.SetReturnReceipts(fakeReceipts)
		fakeUsedGas := uint64(1234)
		processor.SetReturnUsedGas(fakeUsedGas)
		computer := state_computation.NewStateComputer(chain, db, iteratorFactory, processor, trieFactory, validator)
		stateTrie := state_mocks.NewMockStateTrie()
		fakeStateDB := &state.StateDB{}
		stateTrie.SetStateDb(fakeStateDB)
		trieFactory.SetStateTrie(stateTrie)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).NotTo(HaveOccurred())
		validator.AssertValidateStateCalledWith(currentBlock, parentBlock, fakeStateDB, fakeReceipts, fakeUsedGas)
	})

	It("returns error if validating state fails", func() {
		chain, db, iteratorFactory, processor, trieFactory, validator := getMocks()
		validator.SetReturnErr(test_helpers.FakeError)
		computer := state_computation.NewStateComputer(chain, db, iteratorFactory, processor, trieFactory, validator)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("commits validated state to memory database", func() {
		chain, db, iteratorFactory, processor, trieFactory, validator := getMocks()
		computer := state_computation.NewStateComputer(chain, db, iteratorFactory, processor, trieFactory, validator)
		stateTrie := state_mocks.NewMockStateTrie()
		trieFactory.SetStateTrie(stateTrie)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).NotTo(HaveOccurred())
		stateTrie.AssertCommitCalled()
	})

	It("returns error if committing state fails", func() {
		chain, db, iteratorFactory, processor, trieFactory, validator := getMocks()
		computer := state_computation.NewStateComputer(chain, db, iteratorFactory, processor, trieFactory, validator)
		stateTrie := state_mocks.NewMockStateTrie()
		stateTrie.SetReturnErr(test_helpers.FakeError)
		trieFactory.SetStateTrie(stateTrie)
		currentBlock, parentBlock := getFakeBlocks()

		_, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(test_helpers.FakeError))
	})

	It("returns nodes from memory database", func() {
		chain, db, iteratorFactory, processor, trieFactory, validator := getMocks()
		fakeDB := db.CreateFakeUnderlyingDatabase()
		db.SetReturnDatabase(fakeDB)
		iterator := state_mocks.NewMockStateIterator(2)
		iterator.SetReturnHash(common.HexToHash("0x123"))
		iteratorFactory.SetReturnIterator(iterator)
		computer := state_computation.NewStateComputer(chain, db, iteratorFactory, processor, trieFactory, validator)

		currentBlock, parentBlock := getFakeBlocks()

		results, err := computer.ComputeBlockStateTrie(currentBlock, parentBlock)

		Expect(err).NotTo(HaveOccurred())
		Expect(len(results)).To(Equal(2))
		Expect(results[0]).To(Equal([]byte{1, 2, 3, 4, 5}))
	})
})

func getMocks() (*test_helpers.MockBlockChain, *state_mocks.MockStateDatabase, *state_mocks.MockStateIteratorFactory, *test_helpers.MockProcessor, *state_mocks.MockStateTrieFactory, *test_helpers.MockValidator) {
	chain := test_helpers.NewMockBlockChain()
	db := state_mocks.NewMockStateDatabase()
	iteratorFactory := state_mocks.NewMockStateIteratorFactory()
	iterator := state_mocks.NewMockStateIterator(1)
	iteratorFactory.SetReturnIterator(iterator)
	processor := test_helpers.NewMockProcessor()
	trieFactory := state_mocks.NewMockStateTrieFactory()
	stateTrie := state_mocks.NewMockStateTrie()
	trieFactory.SetStateTrie(stateTrie)
	validator := test_helpers.NewMockValidator()
	return chain, db, iteratorFactory, processor, trieFactory, validator
}

func getFakeBlocks() (*types.Block, *types.Block) {
	currentBlockHeader := &types.Header{
		Root:   common.HexToHash("0x123"),
		Number: big.NewInt(456),
	}
	currentBlock := types.NewBlockWithHeader(currentBlockHeader)
	parentBlockHeader := &types.Header{
		Root:   common.HexToHash("0x789"),
		Number: big.NewInt(457),
	}
	parentBlock := types.NewBlockWithHeader(parentBlockHeader)
	return currentBlock, parentBlock
}
