package rawdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

type IAccessorsChain interface {
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetBlockReceipts(hash common.Hash, number uint64) types.Receipts
	GetBody(hash common.Hash, number uint64) *types.Body
	GetCanonicalHash(number uint64) common.Hash
	GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue
}

type AccessorsChain struct {
	ethDbConnection ethdb.Database
}

func NewAccessorsChain(databaseConnection ethdb.Database) *AccessorsChain {
	return &AccessorsChain{ethDbConnection: databaseConnection}
}

func (ldbr *AccessorsChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return rawdb.ReadBlock(ldbr.ethDbConnection, hash, number)
}

func (ldbr *AccessorsChain) GetBlockReceipts(hash common.Hash, number uint64) types.Receipts {
	return rawdb.ReadReceipts(ldbr.ethDbConnection, hash, number)
}

func (ldbr *AccessorsChain) GetBody(hash common.Hash, number uint64) *types.Body {
	return rawdb.ReadBody(ldbr.ethDbConnection, hash, number)
}

func (ldbr *AccessorsChain) GetCanonicalHash(number uint64) common.Hash {
	return rawdb.ReadCanonicalHash(ldbr.ethDbConnection, number)
}

func (ldbr *AccessorsChain) GetHeaderRLP(hash common.Hash, number uint64) rlp.RawValue {
	return rawdb.ReadHeaderRLP(ldbr.ethDbConnection, hash, number)
}
