package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Database struct {
	reader        Reader
	stateComputer IStateComputer
}

func NewLevelDatabase(ldbReader Reader, stateComputer IStateComputer) *Database {
	return &Database{
		reader:        ldbReader,
		stateComputer: stateComputer,
	}
}

func (l Database) ComputeBlockStateTrie(currentBlock *types.Block, parentBlock *types.Block) ([][]byte, error) {
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

func (l Database) GetStateTrieNodes(root []byte) ([][]byte, error) {
	h := common.BytesToHash(root)
	return l.reader.GetStateTrieNodes(h)
}
