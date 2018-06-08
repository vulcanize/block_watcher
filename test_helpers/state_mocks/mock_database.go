package state_mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/vulcanize/block_watcher/pkg/db/level/state_computation"
	"github.com/vulcanize/block_watcher/test_helpers"
)

type MockStateDatabase struct {
	returnDB state.Database
	trie     state_computation.ITrie
}

func NewMockStateDatabase() *MockStateDatabase {
	return &MockStateDatabase{
		returnDB: nil,
		trie:     nil,
	}
}

func (msdb *MockStateDatabase) CreateFakeUnderlyingDatabase() state.Database {
	return &mockStateDatabase{}
}

func (msdb *MockStateDatabase) SetReturnDatabase(db state.Database) {
	msdb.returnDB = db
}

func (msdb *MockStateDatabase) SetReturnTrie(trie state_computation.ITrie) {
	msdb.trie = trie
}

func (msdb *MockStateDatabase) Database() state.Database {
	return msdb.returnDB
}

func (msdb *MockStateDatabase) OpenTrie(root common.Hash) (state_computation.ITrie, error) {
	return msdb.trie, nil
}

func (msdb *MockStateDatabase) TrieDB() state_computation.ITrieDatabase {
	return msdb.returnDB.TrieDB()
}

// implements state.Database interface for testing
type mockStateDatabase struct {
}

func (*mockStateDatabase) ContractCode(addrHash, codeHash common.Hash) ([]byte, error) {
	panic("implement me")
}

func (*mockStateDatabase) ContractCodeSize(addrHash, codeHash common.Hash) (int, error) {
	panic("implement me")
}

func (*mockStateDatabase) CopyTrie(state.Trie) state.Trie {
	panic("implement me")
}

func (*mockStateDatabase) OpenStorageTrie(addrHash, root common.Hash) (state.Trie, error) {
	panic("implement me")
}

func (*mockStateDatabase) OpenTrie(root common.Hash) (state.Trie, error) {
	return &trie.SecureTrie{}, nil
}

func (*mockStateDatabase) TrieDB() *trie.Database {
	trieDB := trie.NewDatabase(&mockEthDB{})
	trieDB.Insert(test_helpers.FakeHash, []byte{1, 2, 3, 4, 5})
	return trieDB
}

// implements eth.Database interface for testing
type mockEthDB struct {
}

func (mockEthDB) Put(key []byte, value []byte) error {
	panic("implement me")
}

func (mockEthDB) Get(key []byte) ([]byte, error) {
	panic("implement me")
}

func (mockEthDB) Has(key []byte) (bool, error) {
	panic("implement me")
}

func (mockEthDB) Delete(key []byte) error {
	panic("implement me")
}

func (mockEthDB) Close() {
	panic("implement me")
}

func (mockEthDB) NewBatch() ethdb.Batch {
	panic("implement me")
}
