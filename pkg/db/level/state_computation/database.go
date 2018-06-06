package state_computation

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

type Database interface {
	Database() state.Database
	TrieDB() *trie.Database
}

type StateDatabase struct {
	db state.Database
}

func NewDatabase(databaseConnection ethdb.Database) *StateDatabase {
	db := state.NewDatabase(databaseConnection)
	return &StateDatabase{db: db}
}

func (db StateDatabase) Database() state.Database {
	return db.db
}

func (db StateDatabase) TrieDB() *trie.Database {
	return db.db.TrieDB()
}
