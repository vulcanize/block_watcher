package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

// Wraps go-ethereum db operations
type Reader interface {
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue
	GetCanonicalHash(number uint64) common.Hash
	GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue
	GetStateTrieNodes(root common.Hash) ([][]byte, error)
}

type LDBReader struct {
	ethDbConnection ethdb.Database
	trieDb          *trie.Database
}

func NewLevelDatabaseReader(databaseConnection ethdb.Database) *LDBReader {
	trieDb := trie.NewDatabase(databaseConnection)
	return &LDBReader{ethDbConnection: databaseConnection, trieDb: trieDb}
}

func (ldbr *LDBReader) GetBlock(hash common.Hash, number uint64) *types.Block {
	return rawdb.ReadBlock(ldbr.ethDbConnection, hash, number)
}

func (ldbr *LDBReader) GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue {
	return rawdb.ReadBodyRLP(ldbr.ethDbConnection, hash, number)
}

func (ldbr *LDBReader) GetCanonicalHash(number uint64) common.Hash {
	return rawdb.ReadCanonicalHash(ldbr.ethDbConnection, number)
}

func (ldbr *LDBReader) GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	return rawdb.ReadHeaderRLP(ldbr.ethDbConnection, hash, number)
}

func (ldbr *LDBReader) GetStateTrieNodes(root common.Hash) ([][]byte, error) {
	var results [][]byte
	stateTrie, err := trie.New(root, ldbr.trieDb)
	if err != nil {
		return nil, err
	}
	stateTrieIterator := stateTrie.NodeIterator(root.Bytes())
	rootNode, err := ldbr.trieDb.Node(root)
	if err != nil {
		return nil, err
	}
	results = append(results, rootNode)
	for stateTrieIterator.Next(true) {
		nextNodeHash := stateTrieIterator.Hash()
		nextNode, err := ldbr.trieDb.Node(nextNodeHash)
		if err != nil {
			return nil, err
		}
		results = append(results, nextNode)
	}
	return results, nil
}
