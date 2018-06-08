package state_computation

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
)

type Database interface {
	Database() state.Database
	OpenTrie(root common.Hash) (ITrie, error)
	TrieDB() ITrieDatabase
}

type StateDatabase struct {
	db     state.Database
	trieDB ITrieDatabase
}

func NewDatabase(databaseConnection ethdb.Database) *StateDatabase {
	db := state.NewDatabase(databaseConnection)
	trieDB := NewTrieDatabase(db.TrieDB())
	return &StateDatabase{db: db, trieDB: trieDB}
}

func (db StateDatabase) Database() state.Database {
	return db.db
}

func (db StateDatabase) OpenTrie(root common.Hash) (ITrie, error) {
	stateTrie, err := db.db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	return NewTrie(stateTrie), nil
}

func (db StateDatabase) TrieDB() ITrieDatabase {
	return db.trieDB
}
