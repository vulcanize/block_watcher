package level

import "github.com/ethereum/go-ethereum/common"

type Database struct {
	reader Reader
}

func NewLevelDatabase(ldbReader Reader) *Database {
	return &Database{
		reader: ldbReader,
	}
}

func (l Database) GetBlockBodyByBlockNumber(blockNumber int64) ([]byte, error) {
	n := uint64(blockNumber)
	h := l.reader.GetCanonicalHash(n)
	return l.reader.GetBodyRLP(h, n), nil
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
