package level

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

// Wraps go-ethereum db operations
type Reader interface {
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetBodyRLP(hash common.Hash, number uint64) rlp.RawValue
	GetCanonicalHash(number uint64) common.Hash
	GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue
}

type LDBReader struct {
	ethDbConnection ethdb.Database
}

func NewLevelDatabaseReader(databaseConnection ethdb.Database) *LDBReader {
	return &LDBReader{ethDbConnection: databaseConnection}
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
