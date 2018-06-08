package state_computation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
)

type ITrieDatabase interface {
	Node(hash common.Hash) ([]byte, error)
}

type TrieDatabase struct {
	db *trie.Database
}

func NewTrieDatabase(db *trie.Database) *TrieDatabase {
	return &TrieDatabase{db: db}
}

func (td *TrieDatabase) Node(hash common.Hash) ([]byte, error) {
	return td.db.Node(hash)
}
