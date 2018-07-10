package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/eth-block-extractor/pkg/wrappers/core/rawdb"
)

type Database struct {
	accessorsChain  rawdb.IAccessorsChain
	stateComputer   IStateComputer
	stateTrieReader IStateTrieReader
}

func NewLevelDatabase(accessorsChain rawdb.IAccessorsChain, stateComputer IStateComputer, stateTrieReader IStateTrieReader) *Database {
	return &Database{
		accessorsChain:  accessorsChain,
		stateComputer:   stateComputer,
		stateTrieReader: stateTrieReader,
	}
}

func (db Database) ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) (common.Hash, error) {
	return db.stateComputer.ComputeBlockStateTrie(currentBlock, parentBlock)
}

func (db Database) GetBlockBodyByBlockNumber(blockNumber int64) ([]byte, error) {
	n := uint64(blockNumber)
	h := db.accessorsChain.GetCanonicalHash(n)
	return db.accessorsChain.GetBodyRLP(h, n), nil
}

func (db Database) GetBlockByBlockNumber(blockNumber int64) *types.Block {
	n := uint64(blockNumber)
	h := db.accessorsChain.GetCanonicalHash(n)
	return db.accessorsChain.GetBlock(h, n)
}

func (db Database) GetBlockHeaderByBlockNumber(blockNumber int64) ([]byte, error) {
	n := uint64(blockNumber)
	h := db.accessorsChain.GetCanonicalHash(n)
	return db.accessorsChain.GetHeaderRLP(h, n), nil
}

func (db Database) GetBlockReceipts(blockNumber int64) types.Receipts {
	n := uint64(blockNumber)
	h := db.accessorsChain.GetCanonicalHash(n)
	return db.accessorsChain.GetBlockReceipts(h, n)
}

func (db Database) GetStateAndStorageTrieNodes(root common.Hash) (stateTrieNodes, storageTrieNodes [][]byte, err error) {
	return db.stateTrieReader.GetStateAndStorageTrieNodes(root)
}
