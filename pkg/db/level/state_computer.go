package level

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

type StateComputer interface {
	ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) ([][]byte, error)
}

type LDBStateComputer struct {
	blockChain      *core.BlockChain
	stateDatabase   state.Database
	stateProcessor  *core.StateProcessor
	stateValidator  core.Validator
	ethDbConnection ethdb.Database
}

func NewLDBStateComputer(databaseConnection ethdb.Database) *LDBStateComputer {
	blockChain, err := core.NewBlockChain(databaseConnection, nil, params.MainnetChainConfig, ethash.NewFaker(), vm.Config{})
	if err != nil {
		panic("Error creating blockchain for state computer")
	}
	stateProcessor := core.NewStateProcessor(params.MainnetChainConfig, blockChain, ethash.NewFaker())
	validator := core.NewBlockValidator(blockChain.Config(), blockChain, blockChain.Engine())
	stateDb := state.NewDatabase(databaseConnection)
	return &LDBStateComputer{
		blockChain:      blockChain,
		stateDatabase:   stateDb,
		stateProcessor:  stateProcessor,
		stateValidator:  validator,
		ethDbConnection: databaseConnection,
	}
}

func (ldbsc *LDBStateComputer) ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) ([][]byte, error) {
	// state.NewDatabase returns a state.Database: backing datastore for state
	// state.New returns a state.StateDB: a new state from a given trie
	stateDb, err := state.New(parentBlock.Root(), ldbsc.stateDatabase)
	if err != nil {
		return nil, fmt.Errorf("Error creating state for block %d: %s\n", parentBlock.Number(), err)
	}

	err = ldbsc.getRootFor(currentBlock, parentBlock, stateDb)
	if err != nil {
		return nil, fmt.Errorf("Error computing root for block %d: %s\n", currentBlock.Number(), err)
	}

	return ldbsc.getNodesFor(stateDb)
}

func (ldbsc *LDBStateComputer) getRootFor(currentBlock *types.Block, parentBlock *types.Block, stateDb *state.StateDB) error {
	receipts, _, usedGas, err := ldbsc.stateProcessor.Process(currentBlock, stateDb, vm.Config{})
	if err != nil {
		return fmt.Errorf("Error processing state: %s\n", err)
	}
	err = ldbsc.stateValidator.ValidateState(currentBlock, parentBlock, stateDb, receipts, usedGas)
	if err != nil {
		return fmt.Errorf("Error validating state: %s\n", err)
	}
	_, err = stateDb.Commit(ldbsc.blockChain.Config().IsEIP158(currentBlock.Number()))
	if err != nil {
		return fmt.Errorf("Error committing state: %s\n", err)
	}
	return nil
}

func (ldbsc *LDBStateComputer) getNodesFor(stateDb *state.StateDB) ([][]byte, error) {
	var results [][]byte
	stateNodeIterator := state.NewNodeIterator(stateDb)
	for stateNodeIterator.Next() {
		stateNodeHash := stateNodeIterator.Hash
		// state trie nodes are sometimes null/empty
		if common.EmptyHash(stateNodeHash) {
			continue
		}
		stateNode, err := ldbsc.stateDatabase.TrieDB().Node(stateNodeHash)
		if err != nil {
			return nil, fmt.Errorf("Error fetching node for hash %s: %s\n", stateNodeHash.String(), err)
		}
		results = append(results, stateNode)
	}
	return results, nil
}
