package state_computation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Computer interface {
	ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) ([][]byte, error)
}

type StateComputer struct {
	blockChain      BlockChain
	db              Database
	iteratorFactory IteratorFactory
	processor       Processor
	trieFactory     TrieFactory
	validator       Validator
}

func NewStateComputer(blockChain BlockChain, db Database, iteratorFactory IteratorFactory, processor Processor, trieFactory TrieFactory, validator Validator) *StateComputer {
	return &StateComputer{
		blockChain:      blockChain,
		db:              db,
		iteratorFactory: iteratorFactory,
		processor:       processor,
		trieFactory:     trieFactory,
		validator:       validator,
	}
}

func (sc *StateComputer) ComputeBlockStateTrie(block, parent *types.Block) ([][]byte, error) {
	stateTrie, err := sc.trieFactory.NewStateTrie(parent.Root(), sc.db.Database())
	if err != nil {
		return nil, err
	}

	err = sc.createStateTrieForBlock(block, parent, stateTrie)
	if err != nil {
		return nil, err
	}

	return sc.persistStateTrieNodes(stateTrie)
}

func (sc *StateComputer) createStateTrieForBlock(block, parent *types.Block, stateTrie Trie) error {
	receipts, _, usedGas, err := sc.processor.Process(block, stateTrie.StateDb())
	if err != nil {
		return err
	}
	err = sc.validator.ValidateState(block, parent, stateTrie.StateDb(), receipts, usedGas)
	if err != nil {
		return err
	}
	_, err = stateTrie.Commit(sc.blockChain.Config().IsEIP158(block.Number()))
	if err != nil {
		return err
	}
	return nil
}

func (sc *StateComputer) persistStateTrieNodes(stateTrie Trie) ([][]byte, error) {
	var results [][]byte
	iterator := sc.iteratorFactory.NewNodeIterator(stateTrie.StateDb())
	for iterator.Next() {
		hash := iterator.Hash()
		// state trie nodes are sometimes null/empty
		if common.EmptyHash(hash) {
			continue
		}
		node, err := sc.db.TrieDB().Node(hash)
		if err != nil {
			return nil, err
		}
		results = append(results, node)
	}
	return results, nil
}
