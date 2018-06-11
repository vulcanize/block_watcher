package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Database struct {
	reader          Reader
	stateComputer   IStateComputer
	stateTrieReader IStateTrieReader
}

func NewLevelDatabase(ldbReader Reader, stateComputer IStateComputer, stateTrieReader IStateTrieReader) *Database {
	return &Database{
		reader:          ldbReader,
		stateComputer:   stateComputer,
		stateTrieReader: stateTrieReader,
	}
}

func (l Database) ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) (common.Hash, error) {
	return l.stateComputer.ComputeBlockStateTrie(currentBlock, parentBlock)
}

func (l Database) GetBlockBodyByBlockNumber(blockNumber int64) ([]byte, error) {
	n := uint64(blockNumber)
	h := l.reader.GetCanonicalHash(n)
	return l.reader.GetBodyRLP(h, n), nil
}

func (l Database) GetBlockByBlockNumber(blockNumber int64) *types.Block {
	n := uint64(blockNumber)
	h := l.reader.GetCanonicalHash(n)
	return l.reader.GetBlock(h, n)
}

func (l Database) GetBlockHeaderByBlockNumber(blockNumber int64) ([]byte, error) {
	n := uint64(blockNumber)
	h := l.reader.GetCanonicalHash(n)
	return l.reader.GetHeaderRLP(h, n), nil
}

func (l Database) GetStateAndStorageTrieNodes(root common.Hash) (stateTrieNodes, storageTrieNodes [][]byte, err error) {
	return l.stateTrieReader.GetStateAndStorageTrieNodes(root)
}
