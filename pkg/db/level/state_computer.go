package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/core"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/core/state"
)

type IStateComputer interface {
	ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) ([][]byte, error)
}

type StateComputer struct {
	blockChain     core.GethCoreBlockChain
	db             state.GethStateDatabase
	processor      core.GethStateProcessor
	stateDBFactory state.GethStateDBFactory
	validator      core.GethBlockValidator
}

func NewStateComputer(blockChain core.GethCoreBlockChain, db state.GethStateDatabase, processor core.GethStateProcessor, stateDBFactory state.GethStateDBFactory, validator core.GethBlockValidator) *StateComputer {
	return &StateComputer{
		blockChain:     blockChain,
		db:             db,
		processor:      processor,
		stateDBFactory: stateDBFactory,
		validator:      validator,
	}
}

func (sc *StateComputer) ComputeBlockStateTrie(block, parent *types.Block) ([][]byte, error) {
	stateTrie, err := sc.stateDBFactory.NewStateDB(parent.Root(), sc.db.Database())
	if err != nil {
		return nil, err
	}

	root, err := sc.createStateTrieForBlock(block, parent, stateTrie)
	if err != nil {
		return nil, err
	}

	return sc.persistStateTrieNodes(root)
}

func (sc *StateComputer) createStateTrieForBlock(block, parent *types.Block, stateTrie state.GethStateDB) (common.Hash, error) {
	var root common.Hash
	receipts, _, usedGas, err := sc.processor.Process(block, stateTrie.StateDB())
	if err != nil {
		return root, err
	}
	err = sc.validator.ValidateState(block, parent, stateTrie.StateDB(), receipts, usedGas)
	if err != nil {
		return root, err
	}
	root, err = stateTrie.Commit(sc.blockChain.Config().IsEIP158(block.Number()))
	if err != nil {
		return root, err
	}
	return root, nil
}

func (sc *StateComputer) persistStateTrieNodes(root common.Hash) ([][]byte, error) {
	var results [][]byte
	stateTrie, err := sc.db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	stateTrieIterator := stateTrie.NodeIterator(nil)
	for stateTrieIterator.Next(true) {
		if stateTrieIterator.Leaf() {
			results = append(results, stateTrieIterator.LeafBlob())
		} else {
			nodeKey := stateTrieIterator.Hash()
			node, err := sc.db.TrieDB().Node(nodeKey)
			if err != nil {
				return nil, err
			}
			results = append(results, node)
		}
	}
	return results, nil
}
